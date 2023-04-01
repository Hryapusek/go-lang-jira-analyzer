package main

import (
	"JiraConnector/apiServer"
)

func main() {
	apiServer.NewServer().Start()
}
