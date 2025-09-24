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

func CreatePost(c *gin.Context) {
	var input struct {
		Content string `form:"content" json:"content"`
	}

	// ambil content dari body
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid := utils.InterfaceToUint(userID)

	// handle upload file
	file, _ := c.FormFile("image")
	var url string
	var publicID string
	if file != nil {
		tempPath := "./temp/" + file.Filename
		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		uploadedURL, uploadedPublicID, err := utils.UploadImage(tempPath, "posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		url = uploadedURL
		publicID = uploadedPublicID
		os.Remove(tempPath)
	}

	// buat post baru
	post := models.Post{
		UserID:   uid,
		Content:  input.Content,
		Image:    &url,
		PublicID: publicID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}
	var newPost models.Post
	if err := database.DB.Preload("User").First(&newPost, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created post"})
		return
	}
	ws.SendToClients(gin.H{"event": "post_created", "data": newPost})
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data":    post,
	})
}

func GetPost(c *gin.Context) {
	var posts []models.Post
	if err := database.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
func GetPostId(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	content := c.PostForm("content")

	updates := map[string]interface{}{
		"content": content,
	}

	file, _ := c.FormFile("image")
	if file != nil {
		if post.PublicID != "" {
			utils.DeleteFromCloudinary(post.PublicID)
		}

		filePath := "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err == nil {
			url, publicID, _ := utils.UploadImage(filePath, "posts")
			updates["image"] = url
			updates["public_id"] = publicID
		}
	}

	if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	var newPost models.Post
	if err := database.DB.Preload("User").First(&newPost, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated post"})
		return
	}

	ws.SendToClients(gin.H{"event": "post_updated", "data": newPost})

	c.JSON(http.StatusOK, newPost)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// hapus gambar di Cloudinary
	if post.PublicID != "" {
		utils.DeleteFromCloudinary(post.PublicID)
	}
	ws.SendToClients(gin.H{"event": "post_deleted", "data": post})
	database.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
