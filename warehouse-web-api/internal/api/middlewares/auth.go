package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"web-api/internal/pkg/auth"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type erpPermissionCacheEntry struct {
	allowed   bool
	expiresAt time.Time
}

var erpPermissionCache sync.Map

func RequireJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "Vui lòng đăng nhập để tiếp tục")
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			return auth.JWTSecret(), nil
		})
		if err != nil || token == nil || !token.Valid {
			response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "Phiên đăng nhập đã hết hạn")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		userID, _ := claims["userid"].(string)
		if strings.TrimSpace(userID) == "" {
			userID, _ = claims["USERID"].(string)
		}
		userID = strings.TrimSpace(userID)
		if !ok || userID == "" {
			response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "Thông tin đăng nhập không hợp lệ")
			ctx.Abort()
			return
		}

		ctx.Set(auth.UserIDContextKey, userID)
		ctx.Next()
	}
}

func RequireERPFormPermission(gsbh, mkid, formID string) gin.HandlerFunc {
	return RequireERPFormPermissionForFeature(gsbh, mkid, formID, "Scan In KTP")
}

func RequireERPFormPermissionForFeature(gsbh, mkid, formID, feature string) gin.HandlerFunc {
	return RequireERPAnyFormPermissionForFeature(gsbh, mkid, feature, formID)
}

func RequireERPAnyFormPermission(gsbh, mkid string, formIDs ...string) gin.HandlerFunc {
	return RequireERPAnyFormPermissionForFeature(gsbh, mkid, "Scan In KTP", formIDs...)
}

func RequireERPAnyFormPermissionForFeature(gsbh, mkid, feature string, formIDs ...string) gin.HandlerFunc {
	normalizedFormIDs := make([]string, 0, len(formIDs))
	seenFormIDs := make(map[string]struct{}, len(formIDs))
	for _, formID := range formIDs {
		formID = strings.ToUpper(strings.TrimSpace(formID))
		if formID == "" {
			continue
		}
		if _, exists := seenFormIDs[formID]; exists {
			continue
		}
		seenFormIDs[formID] = struct{}{}
		normalizedFormIDs = append(normalizedFormIDs, formID)
	}

	return func(ctx *gin.Context) {
		userID, exists := ctx.Get(auth.UserIDContextKey)
		if !exists {
			response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "Vui lòng đăng nhập để tiếp tục")
			ctx.Abort()
			return
		}

		userIDText := strings.TrimSpace(fmt.Sprint(userID))
		cacheKey := strings.Join([]string{gsbh, userIDText, mkid, strings.Join(normalizedFormIDs, ",")}, "|")
		if cached, found := erpPermissionCache.Load(cacheKey); found {
			entry := cached.(erpPermissionCacheEntry)
			if time.Now().Before(entry.expiresAt) {
				if !entry.allowed {
					response.FailWithDetailed(ctx, http.StatusForbidden, nil, fmt.Sprintf("Bạn không có quyền sử dụng chức năng %s", feature))
					ctx.Abort()
					return
				}
				ctx.Next()
				return
			}
			erpPermissionCache.Delete(cacheKey)
		}

		db, err := database.LYS_ERP_Connection()
		if err != nil {
			response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, fmt.Sprintf("Không thể kiểm tra quyền %s", feature))
			ctx.Abort()
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			defer sqlDB.Close()
		}

		var permissionCount int64
		if len(normalizedFormIDs) == 0 {
			response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, fmt.Sprintf("Chưa cấu hình quyền %s", feature))
			ctx.Abort()
			return
		}

		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(normalizedFormIDs)), ",")
		queryArgs := make([]interface{}, 0, 3+len(normalizedFormIDs))
		queryArgs = append(queryArgs, gsbh, userIDText, mkid)
		for _, formID := range normalizedFormIDs {
			queryArgs = append(queryArgs, formID)
		}
		err = db.Raw(`
			SELECT COUNT(1)
			FROM dbo.Blimits
			WHERE GSBH = ? AND USERID = ? AND MKID = ? AND FORMID IN (`+placeholders+`)
		`, queryArgs...).Scan(&permissionCount).Error
		if err != nil {
			response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, fmt.Sprintf("Không thể đọc quyền %s", feature))
			ctx.Abort()
			return
		}

		allowed := permissionCount > 0
		erpPermissionCache.Store(cacheKey, erpPermissionCacheEntry{
			allowed:   allowed,
			expiresAt: time.Now().Add(30 * time.Second),
		})
		if !allowed {
			response.FailWithDetailed(ctx, http.StatusForbidden, nil, fmt.Sprintf("Bạn không có quyền sử dụng chức năng %s", feature))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
