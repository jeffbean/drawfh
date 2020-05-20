package client

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	_writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	_pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	_pingPeriod = (_pongWait * 9) / 10

	_maxMessageSize = 512
)

// Player is the individual client sending and receiving the websocket data.
type Player struct {
	conn *websocket.Conn

	pingTicker *time.Ticker

	// Buffered channel of outbound messages.
	send chan []byte
	// inbound data from the player.
	recv chan []byte

	// channel will close when we quit as a player.
	quit chan struct{}
}

func NewPlayer(name string) *Player {
	return &Player{
		pingTicker: time.NewTicker(_pingPeriod),
		send:       make(chan []byte),
		recv:       make(chan []byte),
		quit:       make(chan struct{}, 1),
	}
}

// Stop will send an optional last message and tear down the player.
func (c *Player) Stop() {
	close(c.send)
	c.pingTicker.Stop()
	c.conn.Close()
}

// Send is sending data to the player. in our case this will
// be some json payload that javascript will be using to update
// the drawing canvas.
func (c *Player) Send(data []byte) {
	select {
	case <-c.quit:
	case c.send <- data:
	}
}

func (c *Player) reader() {
	c.conn.SetReadLimit(_maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(_pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(_pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		select {
		case <-c.quit:
			return
		case c.recv <- message:
		}
	}
}

func (c *Player) processIncoming(data []byte) {
	// todo: process incoming messages from the server as a player.
}

func (c *Player) writer() {
	defer func() {
		close(c.quit)
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(_writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-c.pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(_writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
