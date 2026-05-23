package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

// CasbinAuth 檢查登入的 admin 是否擁有 obj:act 權限。
// superadmin 角色直接放行，不經 casbin 查詢。
func CasbinAuth(enforcer *casbin.Enforcer, obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("admin_role") == "superadmin" {
			c.Next()
			return
		}
		username := c.GetString("admin_username")
		ok, err := enforcer.Enforce(username, obj, act)
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "無此操作權限"})
			return
		}
		c.Next()
	}
}
