package dto

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
}

type UpdateMemberStatusReq struct {
	Status int `json:"status"`
}
