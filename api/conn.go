package api

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sichacvah/portable_chat/model"
)

const (
	WRITE_WAIT  = 10 * time.Second
	PONG_WAIT   = 60 * time.Second
	PING_PERIOD = (PONG_WAIT * 9) / 10
	MAX_SIZE    = 512
)

type WebConn struct {
	WebSocket *websocket.Conn
	Send      chan *model.Message
	UserId    string
}

func NewWebConn(ws *websocket.Conn, userId string) *WebConn {
	return &WebConn{Send: make(chan *model.Message, 64), WebSocket: ws, UserId: userId}
}

func (c *WebConn) readPump() {
	defer func() {
		hub.Unregister(c)
		c.WebSocket.Close()
	}()
	c.WebSocket.SetReadLimit(MAX_SIZE)
	c.WebSocket.SetReadDeadline(time.Now().Add(PONG_WAIT))
	for {
		var msg model.Message
		if err := c.WebSocket.ReadJSON(&msg); err != nil {
			return
		} else {
			msg.UserId = c.UserId
			PublishAndForget(&msg)
		}
	}
}

func (c *WebConn) writePump() {
	ticker := time.NewTicker(PING_PERIOD)

	defer func() {
		ticker.Stop()
		c.WebSocket.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.WebSocket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
				c.WebSocket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.WebSocket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := c.WebSocket.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			c.WebSocket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := c.WebSocket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
