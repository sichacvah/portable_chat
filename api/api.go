package api

func InitApi() {
	r := Srv.Router
	InitUser(r)
	InitWebSocket(r)
}
