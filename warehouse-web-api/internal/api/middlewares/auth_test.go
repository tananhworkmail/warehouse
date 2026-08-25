package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"web-api/internal/api/utils"
	"web-api/internal/pkg/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TestRequireJWTAcceptsLegacyAndMergedUserIDClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, claimName := range []string{"userid", "USERID"} {
		t.Run(claimName, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				claimName: "TEST01",
				"exp":     time.Now().Add(time.Hour).Unix(),
			})
			tokenText, err := token.SignedString(auth.JWTSecret())
			if err != nil {
				t.Fatal(err)
			}

			router := gin.New()
			router.GET("/protected", RequireJWT(), func(ctx *gin.Context) {
				userID, _ := ctx.Get(auth.UserIDContextKey)
				ctx.String(http.StatusOK, userID.(string))
			})

			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			request.Header.Set("Authorization", "Bearer "+tokenText)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != http.StatusOK || response.Body.String() != "TEST01" {
				t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
			}
		})
	}
}

func TestRequireJWTAcceptsTokenFromMergedLogin(t *testing.T) {
	tokenText, err := utils.GenerateJWT("TEST02", "VDH", "Test User", "")
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/protected", RequireJWT(), func(ctx *gin.Context) {
		userID, _ := ctx.Get(auth.UserIDContextKey)
		ctx.String(http.StatusOK, userID.(string))
	})
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+tokenText)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != "TEST02" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
}
