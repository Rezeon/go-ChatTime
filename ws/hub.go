package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan interface{})

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	clients[ws] = true
	for {
		var msg interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			delete(clients, ws)
			break
		}
	}
}

func HandleMessage() {
	for {
		msg := <-broadcast
		for clients := range clients {
			clients.WriteJSON(msg)
		}
	}
}

func SendToClients(msg interface{}) {
	broadcast <- msg
}
