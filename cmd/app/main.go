package main

import (
	"github.com/boingram/aristochat-client"
	"fmt"
)

func main() {
	client, err := aristochat.NewClient("ws://localhost:4000/socket/websocket?username=user2", "test")
	err = aristochat.StartUI(client)
	fmt.Println(fmt.Sprintf("%v", err))
}
