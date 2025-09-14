package controllers

import (
	"gotry/database"
	"gotry/models"
	"gotry/utils"
	"gotry/ws"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CreateMessage membuat pesan baru
func CreateMessage(c *gin.Context) {
	var input struct {
		ReceiverID uint   `form:"receiver_id" json:"receiver_id" binding:"required"`
		Content    string `form:"content" json:"content"`
	}

	// bind data dari form atau JSON
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ambil senderID dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid := utils.InterfaceToUint(userID)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	msg := models.Message{
		SenderID:   uid,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
	}

	// Upload image kalau ada
	file, _ := c.FormFile("image")
	if file != nil {
		tempPath := "./temp/" + file.Filename
		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		url, publicID, err := utils.UploadImage(tempPath, "messages")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		msg.Image = &url
		msg.PublicID = publicID
		os.Remove(tempPath)
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create message"})
		return
	}

	// kirim ke websocket
	ws.SendToClients(gin.H{"event": "message_created", "data": msg})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
		"data":    msg,
	})
}

// GetMessages menampilkan semua pesan
func GetMessages(c *gin.Context) {
	var msgs []models.Message
	if err := database.DB.Preload("Sender").Preload("Receiver").Find(&msgs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "messages not found"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

// GetMessageByID menampilkan pesan berdasarkan ID
func GetMessageByID(c *gin.Context) {
	id := c.Param("id")
	var msg models.Message
	if err := database.DB.Preload("Sender").Preload("Receiver").First(&msg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}
	c.JSON(http.StatusOK, msg)
}

// UpdateMessage memperbarui pesan
func UpdateMessage(c *gin.Context) {
	id := c.Param("id")
	var msg models.Message

	if err := database.DB.First(&msg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// update content
	if content := c.PostForm("content"); content != "" {
		msg.Content = content
	}

	// cek kalau ada file baru
	file, _ := c.FormFile("image")
	if file != nil {
		if msg.PublicID != "" {
			utils.DeleteFromCloudinary(msg.PublicID)
		}

		filePath := "./uploads/" + file.Filename
		c.SaveUploadedFile(file, filePath)

		url, publicID, _ := utils.UploadImage(filePath, "messages")
		msg.Image = &url
		msg.PublicID = publicID
		os.Remove(filePath)
	}

	database.DB.Save(&msg)

	ws.SendToClients(gin.H{"event": "message_updated", "data": msg})

	c.JSON(http.StatusOK, msg)
}

// DeleteMessage menghapus pesan
func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	var msg models.Message

	if err := database.DB.First(&msg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	if msg.PublicID != "" {
		utils.DeleteFromCloudinary(msg.PublicID)
	}

	database.DB.Delete(&msg)

	ws.SendToClients(gin.H{"event": "message_deleted", "data": id})

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
