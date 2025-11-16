package utility

import (
	"chatapp/src/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context) *T {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return nil
	}
	return &body
}

func FormatAndSendResponse[T any](c *gin.Context, res dto.Response[T]) {
	if res.ErrorMessage != "" {
		c.JSON(res.Status, gin.H{"errorMessage": res.ErrorMessage})
		return
	}

	if res.Message != "" {
		c.JSON(res.Status, gin.H{"message": res.Message})
		return
	}

	c.JSON(res.Status, gin.H{"data": res.Data})
}
