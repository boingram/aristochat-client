package aristochat

import (
	"github.com/gorilla/websocket"
	"time"
	"fmt"
)

// Client defines a client that communicates with a Phoenix server over a channel/websocket
type Client struct {
	ch    chan *Payload
	conn  *websocket.Conn
	room  string
	topic string
}

// Message defines the structure of the messages sent across the channel
type Message struct {
	Topic   string  `json:"topic"`
	Event   string  `json:"event"`
	Payload Payload `json:"payload"`
	Ref     string  `json:"ref"`
}

// Payload defines the non-channel-specific body of the Message
type Payload struct {
	Status   string `json:"status"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

// NewClient creates a Client that connects to a phoenix server over a given websocket address
func NewClient(address, room string) (*Client, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		ch: make(chan *Payload),
		conn: conn,
		room: room,
	}

	return &client, nil
}

func (c *Client) Listen() error {
	err := c.join(c.room)
	if err != nil {
		return err
	}
	for {
		time.Sleep(time.Millisecond * 100)
		message, err := c.read()
		if err != nil {
			fmt.Println(fmt.Sprintf("error reading from socket: %v", err))
		}
		if message != nil && message.Event == "new_msg" {
			c.ch <- &message.Payload
		}
	}
}

func (c *Client) Channel() chan *Payload {
	return c.ch
}

// Join joins a given chat room via a join Message on the channel
func (c *Client) join(room string) error {
	topic := "rooms:" + room
	m := Message{
		Topic:   topic,
		Event:   "phx_join",
		Payload: Payload{},
	}
	err := c.write(&m)
	if err != nil {
		return err
	}
	c.topic = topic
	go c.sendHeartbeats()
	return nil
}

// Read reads a Message from the channel
func (c *Client) read() (*Message, error) {
	var msg Message
	err := c.conn.ReadJSON(&msg)
	return &msg, err
}

func (c *Client) sendHeartbeats() {
	m := Message{
		Topic: "phoenix",
		Event: "heartbeat",
	}

	t := time.NewTicker(time.Second * 5)
	go func() {
		for range t.C {
			c.write(&m)
		}
	}()
}

func (c *Client) SendMessage(message string) error {
	p := Payload{
		Body: message,
	}
	m := Message{
		Topic:   c.topic,
		Event:   "new_msg",
		Payload: p,
	}
	return c.write(&m)
}

func (c *Client) write(m *Message) error {
	return c.conn.WriteJSON(m)
}
