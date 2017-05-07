package aristochat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	address := "ws://localhost:8755/websocket"

	server := mockServer{
		t: t,
	}
	// start test server
	go startTestWebsocket(&server)

	c, err := NewClient(address, "test")
	if err != nil {
		t.Fatalf("Error creating new client: %v", err)
	}

	go c.Listen()

	time.Sleep(time.Second * 7)

	err = c.SendMessage("test")
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	if !server.joined {
		t.Fatal("Never joined server!")
	}

	if server.heartbeats == 0 {
		t.Error("No heartbeats received!")
	}

	fromCh := <-c.Channel()

	if fromCh.Body != "test" {
		t.Errorf("Message in of 'test' doesn't equal message out of '%v'", fromCh.Body)
	}
}

type mockServer struct {
	t          *testing.T
	joined     bool
	heartbeats int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func startTestWebsocket(server *mockServer) {
	http.HandleFunc("/websocket", server.listenServer)
	go http.ListenAndServe("localhost:8755", nil)
}

func (m *mockServer) listenServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.t.Fatalf("Error upgrading websocket connection: %v", err)
	}

	for {
		var in Message
		err := conn.ReadJSON(&in)
		if err != nil {
			m.t.Errorf("Error reading from socket: %v", err)
		}
		if in.Event == "heartbeat" {
			if in.Topic != "phoenix" {
				m.t.Errorf("Expected heartbeat to have topic 'phoenix', has: %v", in.Topic)
			}
			m.heartbeats++
		} else if in.Event == "phx_join" {
			if in.Topic != "rooms:test" {
				m.t.Errorf("Expected phx_join to have topic 'rooms:test', has: %v", in.Topic)
			}
			m.joined = true
		} else if in.Event == "new_msg" {
			err = conn.WriteJSON(in)
			if err != nil {
				m.t.Errorf("Error sending message back to client: %v", err)
			}
		}
	}
}
