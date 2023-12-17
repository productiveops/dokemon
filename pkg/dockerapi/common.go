package dockerapi

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func discardIncomingMessages(ws *websocket.Conn) {
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
	}
}

func setupPinging(ws *websocket.Conn, connectionClosed *chan bool) (*sync.Mutex) {
	var mu sync.Mutex
	
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	pingTicker := time.NewTicker(pingPeriod)

	go func() {
		for range pingTicker.C {
			mu.Lock()
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				mu.Unlock()
				if err.Error() != "websocket: close sent" {
					log.Debug().Err(err).Msg("Error when sending ping")
				}
				*connectionClosed <- true
				return
			}
			mu.Unlock()
		}
	}()

	return &mu
}