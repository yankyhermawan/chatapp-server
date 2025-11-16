package routes

import (
	"chatapp/src/message"
	"chatapp/src/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteMap(router *gin.Engine, db *gorm.DB) {
	userRouter := router.Group("/user")
	messageRouter := router.Group("/message")

	userRouter.GET("", user.FindUserHandler(db))
	messageRouter.GET("/all", message.GetAllMessagesHandler(db))
}
