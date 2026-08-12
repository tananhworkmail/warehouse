package controllers

import (
	"fmt"
	"net/http"
	"time"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("TYXUAN@123")

type LoginController struct {
	BaseController
}

var Lg = &LoginController{}

func (h *LoginController) Login(c *gin.Context) {
	var requestParams request.LoginLoginRequest
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		fmt.Println("Error binding JSON:", err)
		return
	}

	// Lấy thông tin người dùng từ database
	login, err := services.LoginServiceInstance.Login(requestParams.UserId, requestParams.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		fmt.Println("Login error:", err)
		return
	}

	fmt.Println("Login:", login)

	token, err := generateToken(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token."})
		fmt.Println("Token generation error:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"role":     login.Role,
		"USERID":   login.UserId,
		"USERNAME": login.UserName,
	})
}

func generateToken(login *types.Login) (string, error) {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"userid": login.UserId,
		"role":   login.Role,
		"name":   login.UserName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
