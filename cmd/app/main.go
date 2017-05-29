package main

import (
	"flag"
	"fmt"
	"github.com/boingram/aristochat-client"
)

func main() {
	var usernamePtr, roomPtr string
	flag.StringVar(&usernamePtr, "username", "", "The user name to chat as")
	flag.StringVar(&roomPtr, "room", "main", "The chat room to join")
	flag.Parse()

	addr := fmt.Sprintf("ws://localhost:4000/socket/websocket?username=%s", usernamePtr)

	client, err := aristochat.NewClient(addr, roomPtr)
	err = aristochat.StartUI(client, usernamePtr)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
}
