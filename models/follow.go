package models

import "gorm.io/gorm"

// Follow = representasi "teman" / "follower"
type Follow struct {
	gorm.Model
	FollowerID uint `json:"follower_id"` // siapa yang follow
	FollowedID uint `json:"followed_id"` // siapa yang di-follow

	Follower User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Followed User `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE"`
}
