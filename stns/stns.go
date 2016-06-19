package stns

import (
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
	var m []rest.Middleware
	if verbose {
		m = rest.DefaultProdStack
	} else {
		m = rest.DefaultCommonStack
	}
	m = append(m, &rest.JsonIndentMiddleware{})

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
					os.Exit(1)
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

	// make pid
	if err := createPidFile(s.pidFileName); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer removePidFile(s.pidFileName)
	log.Printf("Start Server pid:%d", os.Getpid())

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.config.Port), s.Handler()))

}

func (s *Stns) Handler() http.Handler {
	server := rest.NewApi()

	server.Use(s.middleware...)

	// using basic auth
	if s.config.User != "" && s.config.Password != "" {
		var basicAuthMiddleware = &rest.AuthBasicMiddleware{
			Realm: "stns",
			Authenticator: func(user string, password string) bool {
				return user == s.config.User && password == s.config.Password
			},
		}

		// exclude health check
		server.Use(&rest.IfMiddleware{
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

	server.SetApp(router)
	return server.MakeHandler()
}

func (s *Stns) HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}
