package api

import (
	"fmt"
	"net/http"
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sichacvah/portable_chat/model"
	"github.com/sichacvah/portable_chat/utils"
)

type Context struct {
	User model.User
	Err  *model.AppError
}

func (c *Context) SetInvalidParam(where string, name string) {
	c.Err = model.NewAppError(where, "Invalid "+name+" parameter", "")
	c.Err.StatusCode = http.StatusBadRequest
}

func tokenAuthWithUser(uuid string, accessToken string) bool {
	authBackend := utils.InitJWTAuthenticationBackend()
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signign method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
	return err == nil && token.Valid && token.Claims["sub"] == uuid
}

func tokenAuth(accessToken string) bool {
	authBackend := utils.InitJWTAuthenticationBackend()
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signign method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
	return err == nil && token.Valid
}

func RequireContext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := Context{}
	context.Set(r, "context", sessionContext)
	next(w, r)
}

func RequireAuth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	queryParams, err := url.ParseQuery(req.URL.RawQuery)
	token := queryParams["token"][0]
	sessionContext := Context{}
	if err != nil || token == "" {
		sessionContext.Err = model.NewAppError("ServeHttp", "Invalid auth data", "")
		rw.WriteHeader(http.StatusUnauthorized)
	} else {
		context.Set(req, "context", sessionContext)
		next(rw, req)
	}
}

func RequireAuthAndUser(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	queryParams, err := url.ParseQuery(req.URL.RawQuery)
	token := queryParams["token"][0]
	sessionContext := Context{}
	if err != nil || token == "" {
		sessionContext.Err = model.NewAppError("ServeHttp", "Invalid auth data", "")
		rw.WriteHeader(http.StatusUnauthorized)
	} else {
		vars := mux.Vars(req)
		uuid := vars["uuid"]
		if tokenAuthWithUser(string(uuid), string(token)) {
			result := <-Srv.Store.User().Get(uuid)
			if result.Err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
			} else {
				sessionContext.User = result.Data.(model.User)
				context.Set(req, "context", sessionContext)
				next(rw, req)
			}
		}
	}
}
