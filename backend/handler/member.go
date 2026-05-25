package handler

import (
	"errors"
	"strconv"
	"time"

	"artale_market/dto"
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
	var req dto.LoginReq
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
	respOK(c, gin.H{
		"token":    tokenStr,
		"id":       member.ID,
		"nickname": member.Nickname,
		"username": member.Username,
		"email":    member.Email,
		"status":   member.Status,
	})
}

func (h *MemberHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	member, err := h.memberSvc.Register(req.Nickname, req.Username, req.Password, req.Email)

	if err != nil {
		respBadRequest(c, err)
		return
	}
	respOK(c, member)
}

func (h *MemberHandler) Logout(c *gin.Context) {
	respOK(c, gin.H{"message": "logged out"})
}

func (h *MemberHandler) Me(c *gin.Context) {
	idRaw, _ := c.Get("member_id")
	id := uint(idRaw.(float64))
	member, err := h.memberSvc.GetMe(id)
	if err != nil {
		respNotFound(c, err)
		return
	}
	resp := gin.H{
		"id":         member.ID,
		"nickname":   member.Nickname,
		"username":   member.Username,
		"email":      member.Email,
		"status":     member.Status,
		"created_at": member.CreatedAt,
	}
	respOK(c, resp)
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
	var req dto.UpdateMemberStatusReq
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
