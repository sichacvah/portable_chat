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
	sr.Handle("/login", negroni.New(
		negroni.HandlerFunc(RequireContext),
		negroni.HandlerFunc(login),
	)).Methods("POST")
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
	)).Methods("POST")
}

func login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	props := model.MapFromJson(r.Body)
	result := <-Srv.Store.User().GetByLogin(props["login"])
	if result.Err != nil {
		sessionContext.SetInvalidParam("get User by Login", "Login = "+props["login"])
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	user := result.Data.(*model.User)
	if user.ComparePassword(props["password"]) {
		w.WriteHeader(http.StatusOK)
		user.SetToken()
		user.Sanitize()
		msg := &model.Message{}
		msg.UserId = user.Id
		msg.Action = model.ACTION_NEW_USER
		msg.Props = make(map[string]string)
		PublishAndForget(msg)
		w.Write([]byte(user.ToJson()))

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

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
	updatedUser.Sanitize()
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
	user.Sanitize()
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

	cr := <-Srv.Store.User().GetCount()
	if cr.Err != nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	usersCount := cr.Data.(int)
	if usersCount <= 0 {
		user.Role = model.USER_ROLE_ADMIN
	} else {
		user.Role = model.USER_ROLE_USER
	}

	result := <-Srv.Store.User().Save(user)
	if result.Err != nil {
		sessionContext.SetInvalidParam("Create User", "user")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
	createdUser := result.Data.(*model.User)
	createdUser.Sanitize()
	createdUser.SetToken()
	w.Write([]byte(createdUser.ToJson()))
}
