package main

import (
	"github.com/sichacvah/portable_chat/api"
	"github.com/sichacvah/portable_chat/utils"
)

func main() {
	utils.Init()
	api.NewServer()
	api.InitApi()
	api.StartServer()
	api.CreateDefaultChannel()
}
