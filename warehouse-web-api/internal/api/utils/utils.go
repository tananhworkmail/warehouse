// utils/utils.go
package utils

import (
	"time"
	"web-api/internal/pkg/auth"

	"github.com/dgrijalva/jwt-go"
)

// Khóa bí mật để ký token
// Hàm tạo JWT token
func GenerateJWT(USERID string, GSBH string, Username string, SN8 string) (string, error) {
	claims := jwt.MapClaims{
		"GSBH":     GSBH,
		"USERID":   USERID,
		"userid":   USERID,
		"USERNAME": Username,
		"SN8":      SN8,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Hạn token 2 giờ
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(auth.JWTSecret())
}
