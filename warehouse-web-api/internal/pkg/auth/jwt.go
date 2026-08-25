package auth

import (
	"os"
	"strings"
	"web-api/internal/pkg/config"
)

const UserIDContextKey = "authenticatedUserID"

const legacyJWTSecret = "TYXUAN@123"

// JWTSecret returns the one signing key shared by every login endpoint and
// protected API. The environment takes precedence over the server config.
func JWTSecret() []byte {
	secret := strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if secret == "" && config.GetConfig() != nil {
		secret = strings.TrimSpace(config.GetConfig().Server.Secret)
	}
	if secret == "" {
		secret = legacyJWTSecret
	}
	return []byte(secret)
}
