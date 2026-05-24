package handler

import (
	"errors"
	"time"

	"artale_market/dto"
	"artale_market/middleware"
	"artale_market/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AdminHandler struct {
	adminSvc service.AdminService
}

func NewAdminHandler(adminSvc service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, errors.New("username and password required"))
		return
	}

	admin, err := h.adminSvc.Authenticate(req.Username, req.Password)
	if err != nil {
		respUnauthorized(c, errors.New("帳號或密碼錯誤"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(middleware.JwtSecret())
	if err != nil {
		respInternal(c, err)
		return
	}

	respOK(c, gin.H{"token": tokenStr, "admin": admin})
}

func (h *AdminHandler) List(c *gin.Context) {
	admins, err := h.adminSvc.ListAdmins()
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, admins)
}

func (h *AdminHandler) Create(c *gin.Context) {
	var req dto.CreateAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	admin, err := h.adminSvc.CreateAdmin(req.Username, req.Password, req.Role)
	if err != nil {
		respBadRequest(c, err)
		return
	}
	respCreated(c, admin)
}

func (h *AdminHandler) Update(c *gin.Context) {
	var req dto.UpdateAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	admin, err := h.adminSvc.UpdateAdmin(parseID(c), req.Username, req.Password, req.Role)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, admin)
}

func (h *AdminHandler) Delete(c *gin.Context) {
	if err := h.adminSvc.DeleteAdmin(parseID(c)); err != nil {
		respBadRequest(c, err)
		return
	}
	respDeleted(c)
}
