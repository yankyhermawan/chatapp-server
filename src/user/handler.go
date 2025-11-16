package user

import (
	"chatapp/src/dto"
	"chatapp/src/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := utility.BindJSON[dto.RegisterUserBody](ctx)

		response := RegisterUser(db, body)
		utility.FormatAndSendResponse(ctx, response)
	}
}

func LoginUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := utility.BindJSON[dto.LoginUserBody](ctx)
		response := LoginUser(db, body)
		utility.FormatAndSendResponse(ctx, response)
	}
}

func FindUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		search := ctx.Query("search")
		userId := ctx.Query("id")
		response := FindUser(db, search, userId)
		utility.FormatAndSendResponse(ctx, response)
	}
}

func AuthMiddlewareHandler(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	splittedToken := strings.Split(header, " ")
	if splittedToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "Invalid token",
		})
	}
	token := splittedToken[1]
	response := AuthMiddleware(token)
	if response.ErrorMessage != "" {
		utility.FormatAndSendResponse(c, response)
		return
	}
	c.Next()
}
