package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
		c.Set("member_id", claims["sub"])
		c.Set("member_username", claims["username"])
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
		c.Set("admin_id", claims["sub"])
		c.Set("admin_username", claims["username"])
		c.Set("admin_role", claims["role"])
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
