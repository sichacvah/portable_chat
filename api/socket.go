package api

import (
	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sichacvah/portable_chat/model"
	"net/http"
)

func InitWebSocket(r *mux.Router) {
	l4g.Debug("Initializing websocket api")
	r.Handle("/websocket", negroni.New(
		//		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(connect),
	)).Methods("GET")
	hub.Start()
}

func connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l4g.Error("websocket connect err: %v", err)
		return
	}

	sessionContext := context.Get(r, "context")

	wc := NewWebConn(ws, sessionContext.User.Id)
	hub.Register(wc)
	go wc.WritePump()
	wc.ReadPump()
}
