package main

import (
	"fmt"
	"github.com/boingram/aristochat-client"
)

func main() {
	client, err := aristochat.NewClient("ws://localhost:4000/socket/websocket?username=user2")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error creating client: %v", err))
		return
	}
	err = client.Join("test")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error joining room: %v", err))
		return
	}

	p1, err := client.Read()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading msg: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("%v", p1))

	err = client.SendMessage("test")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error sending msg: %v", err))
		return
	}

	p2, err := client.Read()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading msg: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("%v", p2))
}
