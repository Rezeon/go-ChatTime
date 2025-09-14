package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	PublicID  string    `json:"public_id"`
	Profile   *string   `json:"profile"`
	Password  string    `json:"-"`
	Posts     []Post    `gorm:"foreignKey:UserID"`
	Followers []Follow  `gorm:"foreignKey:FollowedID"`
	Following []Follow  `gorm:"foreignKey:FollowerID"`
	Messages  []Message `gorm:"foreignKey:SenderID"`
}

type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
