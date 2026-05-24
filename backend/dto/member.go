package dto

// LoginReq 會員登入請求
type LoginReq struct {
	Username string `json:"username" binding:"required"` // 登入帳號
	Password string `json:"password" binding:"required"` // 登入密碼
}

// UpdateMemberStatusReq 更新會員狀態請求
type UpdateMemberStatusReq struct {
	Status int `json:"status"` // 狀態：1=正常 0=封禁
}
