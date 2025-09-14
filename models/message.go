package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SenderID   uint    `json:"sender_id"`
	ReceiverID uint    `json:"receiver_id"`
	Content    string  `json:"content"`
	PublicID   string  `json:"public_id"`
	Image      *string `json:"image_url"`
	Sender     User    `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Receiver   User    `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE"`
}
