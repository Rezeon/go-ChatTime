package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // untuk menyimpan client yang terhubung
var broadcast = make(chan interface{})       // untuk membuat broadcast dan dikirim kek client

var upgrader = websocket.Upgrader{ // mengubah http menjadi websocket/ws
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) { // untuk menerima koneksi baru dari client
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close() // otomatis menutup koneksi jika funsi selesai

	clients[ws] = true // jika berhasil akan disimpan
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
	for { // menunggu pesan baru dan akan dikirim ke client
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil { // jika error koneksi kek client akan di close
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
