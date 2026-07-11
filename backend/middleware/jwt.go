package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys for propagating auth claims into the standard context.
// Huma handlers receive context.Context (not *gin.Context), so claims
// must be set here to be readable in Huma handler functions.
type CtxAdminID struct{}
type CtxAdminUsername struct{}
type CtxAdminRole struct{}
type CtxMemberID struct{}
type CtxMemberUsername struct{}

func MemberJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret(), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		if claims["type"] != "member" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not a member token"})
			return
		}
		memberID := fmt.Sprintf("%v", claims["sub"])
		memberUsername := fmt.Sprintf("%v", claims["username"])
		c.Set("member_id", memberID)
		c.Set("member_username", memberUsername)
		newCtx := context.WithValue(c.Request.Context(), CtxMemberID{}, memberID)
		newCtx = context.WithValue(newCtx, CtxMemberUsername{}, memberUsername)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := jwtSecret()

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		adminID := fmt.Sprintf("%v", claims["sub"])
		adminUsername := fmt.Sprintf("%v", claims["username"])
		adminRole := fmt.Sprintf("%v", claims["role"])
		c.Set("admin_id", adminID)
		c.Set("admin_username", adminUsername)
		c.Set("admin_role", adminRole)
		newCtx := context.WithValue(c.Request.Context(), CtxAdminID{}, adminID)
		newCtx = context.WithValue(newCtx, CtxAdminUsername{}, adminUsername)
		newCtx = context.WithValue(newCtx, CtxAdminRole{}, adminRole)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}

func JwtSecret() []byte {
	return jwtSecret()
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "artale_market_secret_change_me"
	}
	return []byte(s)
}
