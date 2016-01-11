package api

import (
	"net/http"
	"os"

	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sichacvah/portable_chat/store"
)

type Server struct {
	Store  store.Store
	Router *mux.Router
}

var Srv *Server

func NewServer() {
	l4g.Info("Server is initializing...")

	Srv = &Server{}
	Srv.Store = store.NewBoltDBStore()
	Srv.Router = mux.NewRouter()
}

func StartServer() {
	l4g.Info("Starting Server..")

	var handler http.Handler = Srv.Router
	n := negroni.Classic()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	n.Use(c)
	n.UseHandler(handler)
	http.ListenAndServe(string(os.Args[1])+":5000", n)

}
