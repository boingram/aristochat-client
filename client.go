package aristochat

import "github.com/gorilla/websocket"

// Client defines a client that communicates with a Phoenix server over a channel/websocket
type Client struct {
	conn *websocket.Conn
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
func NewClient(address string) (*Client, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		conn: conn,
	}

	return &client, nil
}

// Join joins a given chat room via a join Message on the channel
func (c *Client) Join(room string) error {
	topic := "rooms:" + room
	m := Message {
		Topic:   topic,
		Event:   "phx_join",
		Payload: Payload{},
	}
	err := c.write(&m)
	if err != nil {
		return err
	}
	c.topic = topic
	return nil
}

// Read reads a Message from the channel
func (c *Client) Read() (*Payload, error) {
	var msg Message
	err := c.conn.ReadJSON(&msg)
	return &msg.Payload, err
}

func (c *Client) SendMessage(message string) error {
	p := Payload {
		Body: message,
	}
	m := Message {
		Topic: c.topic,
		Event: "new_msg",
		Payload: p,
	}
	return c.write(&m)
}

func (c *Client) write(m *Message) error {
	return c.conn.WriteJSON(m)
}




