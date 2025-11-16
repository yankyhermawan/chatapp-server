package message

import (
	"chatapp/src/utility"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllMessagesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		targetId := c.Query("targetId")
		userId, _ := strconv.Atoi(id)
		targetUserId, _ := strconv.Atoi(targetId)
		response := GetAllMessages(db, uint(userId), uint(targetUserId))
		utility.FormatAndSendResponse(c, response)
	}
}
