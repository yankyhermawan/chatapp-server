package database

import "time"

type Message struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	SenderId   uint   `gorm:"column:sender_id" json:"sender_id"`
	ReceiverId uint   `gorm:"column:receiver_id" json:"receiver_id"`
	Message    string `gorm:"column:message" json:"message"`
	Sender     User   `gorm:"foreignKey:SenderId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sender"`
	Receiver   User   `gorm:"foreignKey:ReceiverId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"receiver"`

	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}
