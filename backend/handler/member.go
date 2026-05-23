package handler

import (
	"errors"
	"strconv"
	"time"

	"artale_market/middleware"
	"artale_market/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MemberHandler struct {
	memberSvc service.MemberService
}

func NewMemberHandler(memberSvc service.MemberService) *MemberHandler {
	return &MemberHandler{memberSvc: memberSvc}
}

func (h *MemberHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, errors.New("username and password required"))
		return
	}
	member, err := h.memberSvc.Authenticate(req.Username, req.Password)
	if err != nil {
		respUnauthorized(c, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      member.ID,
		"username": member.Username,
		"type":     "member",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(middleware.JwtSecret())
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"token": tokenStr, "member": member})
}

func (h *MemberHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	members, total, err := h.memberSvc.ListMembers(page, pageSize, search)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{
		"data":      members,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *MemberHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	member, err := h.memberSvc.UpdateStatus(parseID(c), req.Status)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, member)
}

func (h *MemberHandler) Delete(c *gin.Context) {
	if err := h.memberSvc.Delete(parseID(c)); err != nil {
		respNotFound(c, err)
		return
	}
	respDeleted(c)
}
