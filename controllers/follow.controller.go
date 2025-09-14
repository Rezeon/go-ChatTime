package controllers

import (
	"gotry/database"
	"gotry/models"
	"gotry/utils"
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
	followerID := utils.InterfaceToUint(userID) // ubah ke uint jika perlu

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

	c.JSON(http.StatusCreated, gin.H{"message": "Followed successfully", "data": follow})
}

func UnfollowUser(c *gin.Context) {
	var follow models.Follow
	followerID := c.Query("follower_id")
	followedID := c.Query("followed_id")

	if err := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relation not found"})
		return
	}

	if err := database.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	var followers []models.Follow

	if err := database.DB.Preload("Follower").Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

func GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	var following []models.Follow

	if err := database.DB.Preload("Followed").Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following})
}
