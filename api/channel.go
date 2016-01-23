package api

import (
	"fmt"
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
		negroni.HandlerFunc(addMember),
	)).Methods("POST")

	sr.Handle("/{id:[A-Za-z0-9]+}/delete", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(deleteMember),
	)).Methods("POST")

	sr.Handle("/{id:[A-Za-z0-9]+}/join", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(join),
	)).Methods("POST")

	sr.Handle("/{id:[A-Za-z0-9]+}/leave", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(leave),
	)).Methods("POST")

	sr.Handle("/create_direct", negroni.New(
		negroni.HandlerFunc(RequireAuthAndUser),
		negroni.HandlerFunc(createDirectChannel),
	)).Methods("POST")

}

func SetNewChannelAdmin(channelId string, userId string, oldAdminId string) error {
	if len(userId) <= 0 || userId == oldAdminId {
		return model.NewAppError("api.SetNewChannelAdmin", "Wrong User id", "")
	}

	om := <-Srv.Store.Channel().GetMember(channelId, oldAdminId)
	if om.Err != nil {
		return om.Err
	}
	oldAdmin := om.Data.(*model.ChannelMember)
	oldAdmin.Role = model.CHANNEL_ROLE_USER

	nm := <-Srv.Store.Channel().GetMember(channelId, userId)
	if nm.Err != nil {
		return nm.Err
	}

	newAdmin := nm.Data.(*model.ChannelMember)
	newAdmin.Role = model.CHANNEL_ROLE_ADMIN

	result1 := <-Srv.Store.Channel().SaveMember(oldAdmin)
	if result1.Err != nil {
		return result1.Err
	}

	result2 := <-Srv.Store.Channel().SaveMember(newAdmin)
	if result2.Err != nil {
		return result2.Err
	}

	return nil
}

func leave(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)

	props := model.MapFromJson(r.Body)

	channelId := string(vars["id"])
	cr := <-Srv.Store.Channel().Get(channelId)
	if cr.Err != nil {
		sessionContext.SetInvalidParam("leave User from Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cm := <-Srv.Store.Channel().GetMember(channelId, sessionContext.User.Id)
	if cm.Err != nil {
		sessionContext.SetInvalidParam("leave User from Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := cm.Data.(*model.ChannelMember)
	if m.Role == model.CHANNEL_ROLE_ADMIN {
		err := SetNewChannelAdmin(channelId, string(props["user_id"]), sessionContext.User.Id)
		if err != nil {
			sessionContext.SetInvalidParam("leave User from Channel", "Channel ID = "+channelId)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	result := <-Srv.Store.Channel().DeleteMember(m)
	if result.Err != nil {
		sessionContext.SetInvalidParam("leave User from Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func join(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)
	channelId := string(vars["id"])

	cr := <-Srv.Store.Channel().Get(channelId)
	if cr.Err != nil {
		sessionContext.SetInvalidParam("join User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	channel := cr.Data.(model.Channel)
	if channel.Type != model.CHANNEL_OPEN {
		sessionContext.SetInvalidParam("join User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	m := &model.ChannelMember{UserId: sessionContext.User.Id, ChannelId: channel.Id}

	cm := <-Srv.Store.Channel().SaveMember(m)
	if cm.Err != nil {
		sessionContext.SetInvalidParam("join User to Channel", "Channel ID = "+channelId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	savedMember := cm.Data.(*model.ChannelMember)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(savedMember.ToJson()))
}

func CreateDefaultChannel() {
	result := <-Srv.Store.Channel().GetCount()
	if result.Data.(int) <= 0 {
		createDefaultChannel()
	}
}

func createDefaultChannel() {
	channel := model.Channel{Name: model.DEFAULT_CHANNEL, Type: model.CHANNEL_OPEN}
	result := <-Srv.Store.Channel().Save(&channel)
	if result.Err != nil {
		panic(result.Err)
	} else {
		fmt.Println(result.Data)
	}
}

func deleteMember(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	vars := mux.Vars(r)
	channelId := string(vars["id"])
	m := model.ChannelMemberFromJson(r.Body)

	if m.UserId == sessionContext.User.Id {
		sessionContext.SetInvalidParam("delete User from Channel", "cant delete themself")
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	cm := <-Srv.Store.Channel().GetMember(channelId, sessionContext.User.Id)
	if cm.Err != nil {
		sessionContext.SetInvalidParam("delete User from Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currenMember := cm.Data.(*model.ChannelMember)
	if currenMember.Role != model.CHANNEL_ROLE_ADMIN {
		sessionContext.SetInvalidParam("delete User from Channel", "")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	result := <-Srv.Store.Channel().DeleteMember(m)
	if result.Err != nil {
		sessionContext.SetInvalidParam("delete User from Channel", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addMember(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
			if channel.Type == model.CHANNEL_OPEN {
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

func createDirectChannel(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionContext := context.Get(r, "context").(Context)
	props := model.MapFromJson(r.Body)

	otherUserId := string(props["user_id"])

	uc := <-Srv.Store.User().Get(otherUserId)

	if uc.Err != nil {
		sessionContext.SetInvalidParam("create DM Channel", "Channel not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m1 := &model.ChannelMember{
		UserId: sessionContext.User.Id,
		Role:   model.CHANNEL_ROLE_ADMIN,
	}

	m2 := &model.ChannelMember{
		UserId: otherUserId,
		Role:   model.CHANNEL_ROLE_USER,
	}

	channel := new(model.Channel)
	channel.Name = model.GetDMNameFromIds(m1.UserId, m2.UserId)

	cc := <-Srv.Store.Channel().SaveDirectChannel(channel, m1, m2)
	if cc.Err != nil {
		sessionContext.SetInvalidParam("create DM Channel", "Channel not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdChannel := cc.Data.(*model.Channel)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(createdChannel.ToJson()))
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
