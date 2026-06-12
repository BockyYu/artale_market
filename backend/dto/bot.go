package dto

type CreateBotReq struct {
	Name     string `json:"name"     binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=tg line"`
	Token    string `json:"token"    binding:"required"`
	ChatID   string `json:"chat_id"`
}

type UpdateBotReq struct {
	Name     string `json:"name"     binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=tg line"`
	Token    string `json:"token"`
	ChatID   string `json:"chat_id"`
}
