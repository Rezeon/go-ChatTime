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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data":    post,
	})
}

func GetPost(c *gin.Context) {
	var posts []models.Post
	if err := database.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
func GetPostId(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// ambil file baru dari form
	file, _ := c.FormFile("image")
	if file != nil {
		// hapus gambar lama di Cloudinary
		if post.PublicID != "" {
			utils.DeleteFromCloudinary(post.PublicID)
		}

		// upload file baru
		filePath := "./uploads/" + file.Filename
		c.SaveUploadedFile(file, filePath)

		url, publicID, _ := utils.UploadImage(filePath, "posts")

		post.Image = &url
		post.PublicID = publicID
	}

	// update content
	c.Bind(&post)
	database.DB.Save(&post)

	ws.SendToClients(gin.H{"event": "post_created", "data": post})

	c.JSON(http.StatusOK, post)
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

	database.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
