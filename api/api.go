package api

func InitApi() {
	r := Srv.Router
	InitUser(r)
	InitChannel(r)
	InitPost(r)
	InitWebSocket(r)
}
