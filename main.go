package main

import (
	"github.com/gorilla/websocket"
	"fmt"
)

type push struct {
	Topic string `json:"topic"`
	Event string `json:"event"`
	Payload []byte `json:"payload"`
	Ref string `json:"ref"`
}

func writer(conn *websocket.Conn) {
	push := push{
		Topic: "rooms:test",
		Event: "phx_join",
		Payload: []byte{},
		Ref: "",
	}
	conn.WriteJSON(push)
}

func main() {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial("ws://localhost:4000/socket/websocket", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer(conn)
}
