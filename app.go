package main

import (
	"log"
	"net/http"
	"shortlink/middleware"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

//App ...
type App struct {
	Router      *mux.Router
	Middlewares *middleware.Middleware
	Config      *Env
}

//Initialize ...
func (a *App) Initialize(e *Env) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Config = e 
	a.Router = mux.NewRouter()
	a.Middlewares = &middleware.Middleware{}

	a.InitializeRoter()
}

//InitializeRoter ...
func (a *App) InitializeRoter() {
	m := alice.New(a.Middlewares.LoggingHanlder, a.Middlewares.RecoverHanlder)

	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")

}

//Run ...
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
