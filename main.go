package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type event struct {
	Topic   string  `json:"topic"`
	Event   string  `json:"event"`
	Payload payload `json:"payload"`
	Ref     string  `json:"ref"`
}

type payload struct {
	Status   string `json:"status"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

func writer(conn *websocket.Conn) {
	e := event{
		Topic:   "rooms:test",
		Event:   "phx_join",
		Payload: payload{},
	}
	err := conn.WriteJSON(e)
	if err != nil {
		fmt.Println(err)
	}

	msgBody := payload{
		Body: "helloooooooooo",
	}
	e2 := event{
		Topic:   "rooms:test",
		Event:   "new_msg",
		Payload: msgBody,
		Ref:     "",
	}

	err = conn.WriteJSON(e2)
	if err != nil {
		fmt.Println(err)
	}
}

func reader(conn *websocket.Conn) (*event, error) {
	var msg event
	err := conn.ReadJSON(&msg)
	return &msg, err
}

func main() {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial("ws://localhost:4000/socket/websocket?username=user2", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer(conn)
	for {
		time.Sleep(time.Millisecond * 100)
		msg, err := reader(conn)

		if err != nil {
			fmt.Println(fmt.Sprintf("ERROR: %v", err))
		}
		fmt.Println(fmt.Sprintf("%v: %v", msg.Payload.Username, msg.Payload.Body))
	}
}
