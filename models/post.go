package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint    `json:"user_id"`
	User     User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Image    *string `json:"image"`
	PublicID string  `json:"public_id"`
	Content  string  `json:"content"`
}
