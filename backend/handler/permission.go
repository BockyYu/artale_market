package handler

import (
	"artale_market/service"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

// 系統定義的兩個權限
var permDefs = []struct {
	Key string // JSON key
	Obj string // casbin obj
	Act string // casbin act
}{
	{"price_write", "price", "write"},
	{"admin_manage", "admin", "manage"},
}

type PermissionHandler struct {
	enforcer *casbin.Enforcer
	adminSvc service.AdminService
}

func NewPermissionHandler(enforcer *casbin.Enforcer, adminSvc service.AdminService) *PermissionHandler {
	return &PermissionHandler{enforcer: enforcer, adminSvc: adminSvc}
}

// Get 回傳指定管理員的現有權限
func (h *PermissionHandler) Get(c *gin.Context) {
	admin, err := h.adminSvc.GetByID(parseID(c))
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, h.buildPerms(admin.Username))
}

// Update 批次設定/撤銷指定管理員的權限
func (h *PermissionHandler) Update(c *gin.Context) {
	admin, err := h.adminSvc.GetByID(parseID(c))
	if err != nil {
		respNotFound(c, err)
		return
	}

	var req map[string]bool
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}

	for _, p := range permDefs {
		allow, ok := req[p.Key]
		if !ok {
			continue
		}
		if allow {
			h.enforcer.AddPolicy(admin.Username, p.Obj, p.Act)
		} else {
			h.enforcer.RemovePolicy(admin.Username, p.Obj, p.Act)
		}
	}

	respOK(c, h.buildPerms(admin.Username))
}

func (h *PermissionHandler) buildPerms(username string) map[string]bool {
	result := make(map[string]bool, len(permDefs))
	for _, p := range permDefs {
		ok, _ := h.enforcer.Enforce(username, p.Obj, p.Act)
		result[p.Key] = ok
	}
	return result
}
