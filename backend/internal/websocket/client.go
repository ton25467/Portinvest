package websocket

import (
	"context"
	"time"

	"github.com/coder/websocket"
)

const (
	writeWait  = 5 * time.Second
	pingPeriod = 30 * time.Second
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID string
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.Unregister(c)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.Close(websocket.StatusNormalClosure, "channel closed")
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, []byte(`{"type":"ping"}`))
			cancel()
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.Unregister(c)
	}()

	for {
		// Read messages from the client. We discard them as this connection is outbound-broadcast only,
		// but we must keep reading to handle pong/close control frames.
		_, _, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
	}
}

func (c *Client) CloseConn() {
	_ = c.conn.Close(websocket.StatusGoingAway, "closing")
}
