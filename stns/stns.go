package stns

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"
)

type Stns struct {
	config         Config
	configFileName string
	pidFileName    string
	middleware     []rest.Middleware
}

func Create(config Config, configFileName string, pidFileName string, verbose bool) *Stns {
	m := rest.DefaultCommonStack
	m = append(m, &rest.JsonIndentMiddleware{})
	if verbose {
		m = append(m, &rest.AccessLogApacheMiddleware{
			Format: rest.CombinedLogFormat,
		})
	}

	return &Stns{
		config:         config,
		configFileName: configFileName,
		pidFileName:    pidFileName,
		middleware:     m,
	}
}

func (s *Stns) SetMiddleWare(m []rest.Middleware) {
	s.middleware = m
}

func (s *Stns) Start() {
	// wait reload signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM), os.Signal(syscall.SIGUSR1))
	go func() {
		for {
		Loop:
			sig := <-c
			if sig == os.Signal(syscall.SIGUSR1) {
				log.Print("Received signal reload config")
				config, err := LoadConfig(s.configFileName)
				if err != nil {
					log.Fatal(err)
				}
				s.config = config
				log.Printf("Complete reload config\n")
				goto Loop
			} else {
				log.Printf("Received signal '%v', exiting\n", sig)
				removePidFile(s.pidFileName)
				os.Exit(0)
			}
		}
	}()

	server := s.NewHttpServer()

	// make pid
	if err := createPidFile(s.pidFileName); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer removePidFile(s.pidFileName)
	log.Printf("Start Server pid:%d", os.Getpid())

	// tls encryption
	if ok := s.TlsKeysExists(); ok {
		log.Fatal(server.ListenAndServeTLS(s.config.TlsCert, s.config.TlsKey))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func (s *Stns) NewApiHandler() http.Handler {
	api := rest.NewApi()

	api.Use(s.middleware...)

	// using basic auth
	if s.config.User != "" && s.config.Password != "" {
		var basicAuthMiddleware = &rest.AuthBasicMiddleware{
			Realm: "stns",
			Authenticator: func(user string, password string) bool {
				return user == s.config.User && password == s.config.Password
			},
		}

		// exclude health check
		api.Use(&rest.IfMiddleware{
			Condition: func(request *rest.Request) bool {
				return request.URL.Path != "/healthcheck"
			},
			IfTrue: basicAuthMiddleware,
		})
	}

	h := Handler{&s.config}

	router, err := rest.MakeRouter(
		rest.Get("/v2/:resource_name/list", h.GetList),
		rest.Get("/v2/:resource_name/:column/:value", h.Get),
		rest.Get("/:resource_name/list", h.GetList),
		rest.Get("/:resource_name/:column/:value", h.Get),
		rest.Get("/healthcheck", s.HealthChech),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	return api.MakeHandler()
}

func (s *Stns) HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}

func (s *Stns) NewHttpServer() *http.Server {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(s.config.Port),
		Handler: s.NewApiHandler(),
	}

	// tls client authentication
	if _, err := os.Stat(s.config.TlsCa); err == nil {
		ca, err := ioutil.ReadFile(s.config.TlsCa)
		if err != nil {
			log.Fatal(err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(ca)

		tlsConfig := &tls.Config{
			ClientCAs:              caPool,
			SessionTicketsDisabled: true,
			ClientAuth:             tls.RequireAndVerifyClientCert,
		}

		tlsConfig.BuildNameToCertificate()
		server.TLSConfig = tlsConfig
	}

	return server
}

func (s *Stns) TlsKeysExists() bool {
	if s.config.TlsCert != "" && s.config.TlsKey != "" {
		for _, v := range []string{s.config.TlsCert, s.config.TlsKey} {
			if _, err := os.Stat(v); err != nil {
				log.Fatal(err)
			}
		}
		return true
	}
	return false
}
