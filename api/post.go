package api

import (
	"net/http"

	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sichacvah/portable_chat/model"
)

func InitPost(r *mux.Router) {
	l4g.Debug("Initializing post api routes")

	r.Handle("/posts/{post_id}", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(getPost),
	)).Methods("GET")

	sr := r.PathPrefix("/channels/{id:[A-Za-z0-9]+}").Subrouter()

	sr.Handle("/create", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(createPost),
	)).Methods("POST")

	sr.Handle("/update", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(updatePost),
	)).Methods("POST")

	sr.Handle("/posts/{offset:[0-9]+}/{limit:[0-9]+}", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(getPosts),
	)).Methods("POST")
}

func GetPosts(channelId string) (map[string]*model.Post, *model.AppError) {
	result := <-Srv.Store.Post().GetPosts(channelId)
	if result.Err != nil {
		return nil, result.Err
	}

	posts := result.Data.(map[string]*model.Post)
	if posts == nil {
		return nil, model.NewAppError("Get Posts", "Posts not found", "")
	}

	return posts, nil
}

func getPosts(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)

	channelId := string(vars["id"])

	posts, err := GetPosts(channelId)
	if err != nil {
		sessionContext.SetInvalidParam("Error while get posts", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(model.PostsMapToJson(posts)))
}

func getPost(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)

	postId := string(vars["post_id"])

	post, err := GetPost(postId)
	if err != nil {
		sessionContext.SetInvalidParam("Error while get post", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(post.ToJson()))
}

func GetPost(postId string) (*model.Post, *model.AppError) {
	result := <-Srv.Store.Post().Get(postId)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Data.(*model.Post) == nil {
		return nil, model.NewAppError("Get Post", "Post not found", "")
	}

	return result.Data.(*model.Post), nil
}

func updatePost(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	// vars := mux.Vars(r)
	// channelId := string(vars["id"])

	post := model.PostFromJson(r.Body)

	if post == nil || post.Id == "" {
		sessionContext.SetInvalidParam("Invalid post while create", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	rPost, err := UpdatePost(post)
	if err != nil {
		sessionContext.SetInvalidParam("Invalid post while create", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rPost.ToJson()))
}

func UpdatePost(post *model.Post) (*model.Post, *model.AppError) {
	result := <-Srv.Store.Post().Update(post)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Post), nil
}

func CreatePost(post *model.Post) (*model.Post, *model.AppError) {
	result := <-Srv.Store.Post().Save(post)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Post), nil
}

func createPost(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)
	channelId := string(vars["id"])

	post := model.PostFromJson(r.Body)

	if post == nil {
		sessionContext.SetInvalidParam("Invalid post while create", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}

	post.ChannelId = channelId
	post.UserId = sessionContext.User.Id

	rPost, err := CreatePost(post)
	if err != nil {
		sessionContext.SetInvalidParam("Invalid post while create", "")
		w.WriteHeader(sessionContext.Err.StatusCode)
		return
	}
	message := &model.Message{}

	message.UserId = rPost.UserId
	message.ChannelId = rPost.ChannelId
	message.Action = rPost.
		PublishAndForget()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rPost.ToJson()))
}
