package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type Response[T any] struct {
	Status       int    `json:"status"`
	Data         T      `json:"data,omitempty"`
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type RegisterUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	UserId uint   `json:"userId"`
	Token  string `json:"token"`
}

type JwtPayload struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type FindUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type ReceivedMessage struct {
	Message    string `json:"message"`
	ReceiverId uint   `json:"receiver_id"`
	SenderId   uint   `json:"sender_id"`
}

type SendMessage struct {
	Message    string `json:"message"`
	SenderId   uint   `json:"sender_id"`
	ReceiverId uint   `json:"receiver_id"`
}
