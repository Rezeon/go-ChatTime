package controllers

import (
	"gotry/database"
	"gotry/models"
	"gotry/utils"
	"gotry/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FollowUser(c *gin.Context) {
	// Ambil follower_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	followerID := utils.InterfaceToUint(userID)

	var input struct {
		FollowedID uint `form:"followed_id" json:"followed_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah sudah follow
	var existing models.Follow
	if err := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, input.FollowedID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already following"})
		return
	}

	follow := models.Follow{
		FollowerID: followerID,
		FollowedID: input.FollowedID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}

	ws.SendToClients(gin.H{"event": "follow_created", "data": follow})
	c.JSON(http.StatusCreated, gin.H{"message": "Followed successfully", "data": follow})
}

func UnfollowUser(c *gin.Context) {
	var follow models.Follow
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	followerID := utils.InterfaceToUint(userID)
	followedID := c.Query("followed_id")

	if err := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relation not found"})
		return
	}

	if err := database.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
		return
	}
	ws.SendToClients(gin.H{"event": "follow_unfollow", "data": follow})
	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func GetFollowers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var followers []models.Follow

	if err := database.DB.Preload("Follower").Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

func GetFollowing(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var following []models.Follow

	if err := database.DB.Preload("Followed").Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following})
}
