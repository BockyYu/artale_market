package dto

type AdminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateAdminReq 新增管理員請求
type CreateAdminReq struct {
	Username string `json:"username" binding:"required"` // 登入帳號
	Password string `json:"password" binding:"required"` // 登入密碼
	Role     string `json:"role"`                        // 角色，預設 admin
}

// UpdateAdminReq 更新管理員資訊請求（所有欄位選填）
type UpdateAdminReq struct {
	Username string `json:"username"` // 新帳號名稱
	Password string `json:"password"` // 新密碼
	Role     string `json:"role"`     // 新角色
}
