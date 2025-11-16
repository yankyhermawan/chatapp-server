package user

import (
	"chatapp/database"
	"chatapp/src/dto"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, data *dto.RegisterUserBody) dto.Response[database.User] {
	var user database.User

	err := db.First(&user, "username=?", data.Username).Error

	if err == nil {
		return dto.Response[database.User]{
			Status:       http.StatusConflict,
			ErrorMessage: "User already exist",
		}
	}

	saltStr := os.Getenv("HASH_SALT")
	salt, _ := strconv.Atoi(saltStr)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(data.Password), salt)
	user = database.User{
		Username: data.Username,
		Password: string(hashed),
		Name:     data.Name,
	}

	db.Create(&user)

	return dto.Response[database.User]{
		Status:  http.StatusCreated,
		Message: "User Created",
	}
}

func LoginUser(db *gorm.DB, data *dto.LoginUserBody) dto.Response[dto.LoginUserResponse] {
	var user database.User

	err := db.First(&user, "username=?", data.Username).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.Response[dto.LoginUserResponse]{
			Status:       http.StatusNotFound,
			ErrorMessage: "User not Found",
		}
	}
	jwtkey := os.Getenv("JWT_KEY")
	jwtByte := []byte(jwtkey)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return dto.Response[dto.LoginUserResponse]{
			Status:       http.StatusBadRequest,
			ErrorMessage: "Password Mismatch",
		}
	}
	expTime := time.Now().Add(24 * time.Hour)
	payload := dto.JwtPayload{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	tokenString, _ := token.SignedString(jwtByte)

	responseBody := dto.LoginUserResponse{
		Token:  tokenString,
		UserId: user.ID,
	}

	return dto.Response[dto.LoginUserResponse]{
		Status: http.StatusOK,
		Data:   responseBody,
	}
}

func FindUser(db *gorm.DB, search string, id string) dto.Response[[]dto.FindUserResponse] {
	var users []database.User
	var conditions []string
	var args []interface{}
	var err error

	if search != "" {
		searchLike := fmt.Sprintf("%%%s%%", search)
		conditions = append(conditions, "username LIKE ?")
		args = append(args, searchLike)
	}
	if id != "" {
		intId, _ := strconv.Atoi(id)
		conditions = append(conditions, "id = ?")
		args = append(args, uint(intId))
	}
	if len(conditions) > 0 {
		err = db.Limit(10).Where(strings.Join(conditions, " AND "), args...).Find(&users).Error
	} else {
		err = db.Limit(10).Find(&users).Error
	}
	if err != nil {
		return dto.Response[[]dto.FindUserResponse]{
			Status:       http.StatusNotFound,
			ErrorMessage: "User Not Found",
		}
	}
	results := make([]dto.FindUserResponse, 0, len(users))
	for _, u := range users {
		results = append(results, dto.FindUserResponse{
			ID:       u.ID,
			Username: u.Username,
			Name:     u.Name,
		})
	}

	return dto.Response[[]dto.FindUserResponse]{
		Status: http.StatusOK,
		Data:   results,
	}
}

func AuthMiddleware(token string) dto.Response[any] {
	jwtkey := os.Getenv("JWT_KEY")
	jwtByte := []byte(jwtkey)
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtByte, nil
	})

	if err != nil {
		return dto.Response[any]{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "Internal Server Error",
		}
	}

	if !parsed.Valid {
		return dto.Response[any]{
			Status:       http.StatusUnauthorized,
			ErrorMessage: "Invalid token",
		}
	}
	return dto.Response[any]{
		Status:  http.StatusOK,
		Message: "OK",
	}
}
