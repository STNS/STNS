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
	config         *Config
	configFileName string
	pidFileName    string
}

func Create(config *Config, configFileName string, pidFileName string) *Stns {
	return &Stns{config, configFileName, pidFileName}
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
	server.Use(rest.DefaultDevStack...)
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

	router, err := rest.MakeRouter(
		rest.Get("/:resource_name/list", s.GetList),
		rest.Get("/:resource_name/:column/:value", s.Get),
		rest.Get("/healthcheck", s.HealthChech),
	)
	if err != nil {
		log.Fatal(err)
	}

	server.SetApp(router)
	return server.MakeHandler()
}

func (s *Stns) Get(w rest.ResponseWriter, r *rest.Request) {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")
	query := Query{s.config, resource_name, column, value}
	query.Response(w, r)
}

func (s *Stns) GetList(w rest.ResponseWriter, r *rest.Request) {
	resource_name := r.PathParam("resource_name")
	query := Query{s.config, resource_name, "list", ""}
	query.Response(w, r)
}

func (s *Stns) HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}
