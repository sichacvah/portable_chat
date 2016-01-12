package api

import (
	"net/http"

	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func InitWebSocket(r *mux.Router) {
	l4g.Debug("Initializing websocket api")
	r.Handle("/websocket/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(connect),
	)).Methods("GET")
	hub.Start()
}

func connect(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		l4g.Error("websocket connect err: %v", err)
		return
	}

	sessionContext := context.Get(req, "context").(Context)

	wc := NewWebConn(ws, sessionContext.User.Id)
	hub.Register(wc)
	go wc.writePump()
	wc.readPump()
}
