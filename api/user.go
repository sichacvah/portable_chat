package api

import (
	"net/http"

	"github.com/sichacvah/portable_chat/model"

	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//InitUser Initialize user routes
func InitUser(r *mux.Router) {
	l4g.Debug("Initializing user api routes")
	sr := r.PathPrefix("/users").Subrouter()
	sr.Handle("/create", negroni.New(
		negroni.HandlerFunc(RequireContext),
		negroni.HandlerFunc(createUser),
	)).Methods("POST")
	sr.Handle("/", negroni.New(
		negroni.HandlerFunc(RequireTokenAuthentication),
		negroni.HandlerFunc(allUsers),
	)).Methods("GET")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireTokenAuthentication),
		negroni.HandlerFunc(allUsers),
	)).Methods("GET")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireTokenAuthenticationAndUser),
		negroni.HandlerFunc(deleteUser),
	)).Methods("DELETE")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireTokenAuthenticationAndUser),
		negroni.HandlerFunc(updateUser),
	))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := model.UserFromJson(r.Body)
	sessionContext := context.Get(r, "context")
	w.Header().Set("Content-Type", "application/json")
	if user == nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	result := <-Srv.Store.UserStore().Save(user)
	if result.Err != nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.data.ToJson()))
}
