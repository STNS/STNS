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

// Stns stns object
type Stns struct {
	config         Config
	configFileName string
	pidFileName    string
	middleware     []rest.Middleware
}

// NewServer make server object
func NewServer(config Config, configFileName string, pidFileName string, verbose bool) *Stns {
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

// SetMiddleWare use only test
func (s *Stns) SetMiddleWare(m []rest.Middleware) {
	s.middleware = m
}

// Start server start
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

	server := s.newHTTPServer()

	// make pid
	if err := createPidFile(s.pidFileName); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer removePidFile(s.pidFileName)
	log.Printf("Start Server pid:%d", os.Getpid())

	// tls encryption
	if ok := s.tlsKeysExists(); ok {
		log.Fatal(server.ListenAndServeTLS(s.config.TLSCert, s.config.TLSKey))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func (s *Stns) newAPIHandler() http.Handler {
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
		rest.Get("/v3/:resource_name/list", h.GetList),
		rest.Get("/v3/:resource_name/:column/:value", h.Get),
		rest.Get("/v2/:resource_name/list", h.GetList),
		rest.Get("/v2/:resource_name/:column/:value", h.Get),
		rest.Get("/:resource_name/list", h.GetList),
		rest.Get("/:resource_name/:column/:value", h.Get),
		rest.Get("/healthcheck", s.healthCheck),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	return api.MakeHandler()
}

func (s *Stns) healthCheck(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}

func (s *Stns) newHTTPServer() *http.Server {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(s.config.Port),
		Handler: s.newAPIHandler(),
	}

	// tls client authentication
	if _, err := os.Stat(s.config.TLSCa); err == nil {
		ca, err := ioutil.ReadFile(s.config.TLSCa)
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

func (s *Stns) tlsKeysExists() bool {
	if s.config.TLSCert != "" && s.config.TLSKey != "" {
		for _, v := range []string{s.config.TLSCert, s.config.TLSKey} {
			if _, err := os.Stat(v); err != nil {
				log.Fatal(err)
			}
		}
		return true
	}
	return false
}
