package message

import (
	"chatapp/database"
	"chatapp/src/dto"
	"net/http"

	"gorm.io/gorm"
)

func CreateMessage(db *gorm.DB, message dto.ReceivedMessage) {
	messageToBeCreate := database.Message{
		SenderId:   message.SenderId,
		ReceiverId: message.ReceiverId,
		Message:    message.Message,
	}
	db.Create(&messageToBeCreate)
}

func GetAllMessages(db *gorm.DB, userId uint, targetId uint) dto.Response[[]database.Message] {
	var messages []database.Message

	db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, targetId, targetId, userId).
		Find(&messages)
	return dto.Response[[]database.Message]{
		Status: http.StatusOK,
		Data:   messages,
	}
}
