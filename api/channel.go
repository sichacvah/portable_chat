package api

import (
	"net/http"
	"strings"

	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sichacvah/portable_chat/model"
)

func InitChannel(r *mux.Router) {
	l4g.Debug("Initializing channel api routes")

	sr := r.PathPrefix("/channels").Subrouter()
	sr.Handle("/", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(getChannels),
	)).Methods("GET")
	sr.Handle("/create", negroni.New(
		negroni.HandlerFunc(RequireAuth),
		negroni.HandlerFunc(createChannel),
	)).Methods("POST")
	sr.Handle("/update", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(updateChannel),
	)).Methods("POST")
	sr.Handle("/{id:[A-Za-z0-9]+}/add", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(addUser),
	)).Methods("POST")
}

func addUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	newMember := model.ChannelMemberFromJson(r.Body)
	vars := mux.Vars(r)
	channelId := string(vars["id"])

	uc := <-Srv.Store.Channel().GetMember(channelId, sessionContext.User.Id)
	if uc.Err != nil {
		sessionContext.SetInvalidParam("add User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	channelMember := uc.Data.(*model.ChannelMember)
	if channelMember.Role != model.CHANNEL_ROLE_ADMIN {
		sessionContext.SetInvalidParam("add User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mc := <-Srv.Store.Channel().SaveMember(newMember)
	if mc.Err != nil {
		sessionContext.SetInvalidParam("add User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newMember.ToJson()))
}

func updateChannel(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)

	channel := model.ChannelFromJson(r.Body)

	if channel == nil {
		sessionContext.SetInvalidParam("create Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !channel.IsValid() {
		sessionContext.SetInvalidParam("create Channel", "Channel not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := <-Srv.Store.Channel().Save(channel)
	if result.Err != nil {
		sessionContext.SetInvalidParam("create Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	createdChannel := result.Data.(*model.Channel)
	w.Write([]byte(createdChannel.ToJson()))
}

func InChannelMembers(userId string, members map[*model.ChannelMember]bool) bool {
	for channelMember, _ := range members {
		if userId == channelMember.UserId {
			return true
		}
	}
	return false
}

func getChannels(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)

	result := <-Srv.Store.Channel().GetChannels(sessionContext.User.Id)
	if result.Err != nil {
		sessionContext.SetInvalidParam("get Channels", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	visibleChannels := make(map[string]*model.Channel)
	// channelsResult := make(map[string]string)

	channels := result.Data.(map[*model.Channel]bool)
	for channel, ok := range channels {
		if ok {
			if channel.Type != model.CHANNEL_OPEN {
				visibleChannels[channel.Id] = channel
			} else {
				mr := <-Srv.Store.Channel().GetMembers(channel)
				if mr.Err == nil {
					members := mr.Data.(map[*model.ChannelMember]bool)
					if InChannelMembers(sessionContext.User.Id, members) {
						visibleChannels[channel.Id] = channel
					}
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(model.ChannelMapToJson(visibleChannels)))
}

func createChannel(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	channel := model.ChannelFromJson(r.Body)

	if channel == nil {
		sessionContext.SetInvalidParam("create Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !channel.IsValid() {
		sessionContext.SetInvalidParam("create Channel", "Channel not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if channel.Type == model.CHANNEL_DIRECT {
		sessionContext.SetInvalidParam("createDirectChannel", "Must use createDirectChannel api service for direct message channel creation")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.Index(channel.Name, "__") > 0 {
		sessionContext.SetInvalidParam("createDirectChannel", "Invalid character '__' in channel name for non-direct channel")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := <-Srv.Store.Channel().Save(channel)
	if result.Err != nil {
		sessionContext.SetInvalidParam("create Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	firstMember := model.ChannelMember{UserId: sessionContext.User.Id, ChannelId: channel.Id, Role: model.CHANNEL_ROLE_ADMIN}
	rm := <-Srv.Store.Channel().SaveMember(&firstMember)
	if rm.Err != nil {
		sessionContext.SetInvalidParam("create Channel", "create member")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	createdChannel := result.Data.(*model.Channel)
	w.Write([]byte(createdChannel.ToJson()))

}
