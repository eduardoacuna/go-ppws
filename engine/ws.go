package engine

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	// Period allowed to write a message through a websocket
	writeWaitPeriod = 10 * time.Second

	// Period between websocket ping messages
	pingTickPeriod = 60 * time.Second

	// Period allowed to read a message from a websocket
	readWaitPeriod = pingTickPeriod + 5*time.Second
)
