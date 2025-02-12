package main

import (
	"ht-crm/autoconfigure"
	"ht-crm/src/ht/web"
)

func main() {
	autoconfigure.InitConfig("")
	web.Start()
}
