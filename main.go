package main

import (
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/src/ht/web"
)

func main() {
	autoconfigure.InitConfig("")
	web.Start()
}
