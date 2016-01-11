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
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(allUsers),
	)).Methods("GET")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(getUser),
	)).Methods("GET")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(deleteUser),
	)).Methods("DELETE")
	sr.Handle("/{uuid}", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(updateUser),
	))
}

func deleteUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	if sessionContext.Err != nil {
		sessionContext.SetInvalidParam("all Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	vars := mux.Vars(r)
	userId := vars["uuid"]

	result := <-Srv.Store.User().Delete(string(userId))
	if result.Err != nil {
		sessionContext.SetInvalidParam("update Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func updateUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	if sessionContext.Err != nil {
		sessionContext.SetInvalidParam("all Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	user := model.UserFromJson(r.Body)
	result := <-Srv.Store.User().Update(user)
	if result.Err != nil {
		sessionContext.SetInvalidParam("update Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	updatedUser := result.Data.(*model.User)
	w.Write([]byte(updatedUser.ToJson()))
}

func getUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	if sessionContext.Err != nil {
		sessionContext.SetInvalidParam("all Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	vars := mux.Vars(r)
	userId := vars["uuid"]
	result := <-Srv.Store.User().Get(string(userId))
	if result.Err != nil {
		sessionContext.SetInvalidParam("get Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	user := result.Data.(*model.User)
	w.Write([]byte(user.ToJson()))
}

func allUsers(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	w.Header().Set("Content-Type", "application/json")
	result := <-Srv.Store.User().GetUsers()
	if result.Err != nil || sessionContext.Err != nil {
		sessionContext.SetInvalidParam("all Users", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	w.Write(result.Data.([]byte))
}

func createUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := model.UserFromJson(r.Body)
	sessionContext := context.Get(r, "context").(Context)
	w.Header().Set("Content-Type", "application/json")
	if user == nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	result := <-Srv.Store.User().Save(user)
	if result.Err != nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.Data.(*model.User).ToJson()))
}
