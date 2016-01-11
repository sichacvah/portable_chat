package main

import (
	"github.com/sichacvah/portable_chat/store"
	"github.com/sichacvah/portable_chat/utils"
)

func main() {
	utils.Init()
	store.Init()
	api.InitApi()
}
