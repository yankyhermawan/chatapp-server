package database

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"column:username; unique" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}
