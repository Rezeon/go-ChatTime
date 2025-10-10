package ws

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var clients = make(map[uint]*websocket.Conn) // menyimpan koneksi client
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // sementara allow all origin
}

func HandleConnections(c *gin.Context) {
	// Ambil token dari query parameter
	tokenString := c.Query("token")
	if tokenString == "" {
		http.Error(c.Writer, "No token provided", http.StatusUnauthorized)
		return
	}

	// Hapus "Bearer " jika ada (tidak wajib, tapi aman)
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	jwtSecret := []byte(os.Getenv("JWT_TOKEN"))
	if len(jwtSecret) == 0 {
		fmt.Println("JWT_TOKEN is not set in environment")
		http.Error(c.Writer, "Server error", http.StatusInternalServerError)
		return
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		fmt.Println("JWT parse error:", err)
		http.Error(c.Writer, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Ambil userID dari klaim
	var userID uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, ok := claims["user_id"].(float64); ok {
			userID = uint(id)
		} else {
			http.Error(c.Writer, "user_id not found in token", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(c.Writer, "invalid token claims", http.StatusUnauthorized)
		return
	}

	// Upgrade koneksi HTTP ke WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Simpan koneksi user
	clients[userID] = ws
	defer func() {
		fmt.Printf("User %d disconnected\n", userID)
		delete(clients, userID)
	}()

	fmt.Printf("User %d connected\n", userID)

	// Loop membaca pesan dari client
	for {
		var msg struct {
			ToUserID uint        `json:"to"`
			Content  interface{} `json:"msg"`
		}
		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println("Read error:", err)
			break
		}

		// Kirim ke user tujuan
		SendToUser(msg.ToUserID, map[string]interface{}{
			"from": userID,
			"msg":  msg.Content,
		})
	}
}

// Kirim pesan ke user tertentu
func SendToUser(userID uint, msg interface{}) {
	if conn, ok := clients[userID]; ok {
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Write error:", err)
			conn.Close()
			delete(clients, userID)
		}
	} else {
		fmt.Printf("User %d not online\n", userID)
	}
}
