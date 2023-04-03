package main

import (
	"github.com/piyush1146115/shell-mirror/api"
)

func main() {
	api.CreateServer("8088")
	api.StartServer()
}
