package main

import (
	"flag"
	"fmt"
	"github.com/boingram/aristochat-client"
	"os"
	"strings"
)

func main() {
	var usernamePtr, roomPtr, serverPtr string
	flag.StringVar(&usernamePtr, "username", "", "The user name to chat as")
	flag.StringVar(&roomPtr, "room", "main", "The chat room to join")
	flag.StringVar(&serverPtr, "server", "", "The aristochat server to connect to")
	flag.Parse()

	if usernamePtr == "" || roomPtr == "" || serverPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	addr := fmt.Sprintf("ws://%s/socket/websocket?username=%s", strings.Replace(serverPtr, "http://", "", -1), usernamePtr)

	client, err := aristochat.NewClient(addr, roomPtr)
	err = aristochat.StartUI(client, usernamePtr)
	if err != nil {
		fmt.Println("Error starting Aristochat Client!")
	}
}
