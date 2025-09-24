package ws

import (
	"fmt"
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
		fmt.Println("❌ Upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	fmt.Println("WebSocket connection")

	for {
		var msg interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println("error:", err)
			delete(clients, ws)
			break
		}
	}
}

func HandleMessage() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SendToClients(msg interface{}) {
	broadcast <- msg
}
