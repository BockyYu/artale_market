package dto

// ScrollSearchReq 卷軸列表查詢請求
type ScrollSearchReq struct {
	Date       string   `json:"date"`        // 查詢日期，格式 2006-01-02，空值預設今天
	Percentage []int    `json:"percentage"`  // 成功率篩選，空值為全部
	Category   []string `json:"category"`    // 分類篩選，傳入 scroll_all 視為全部卷軸
	SortBy     string   `json:"sort_by"`     // 排序方式：price_desc / price_asc / change_desc / change_asc / percentage_asc / percentage_desc
	Page       int      `json:"page"`        // 頁碼，從 1 開始
	PageSize   int      `json:"page_size"`   // 每頁筆數
}

// SkillBookSearchReq 技能書列表查詢請求
type SkillBookSearchReq struct {
	Date     string   `json:"date"`       // 查詢日期，格式 2006-01-02，空值預設今天
	Category []string `json:"category"`   // 職業篩選，空值為全部職業
	SortBy   string   `json:"sort_by"`    // 排序方式：price_desc / price_asc / change_desc / change_asc / percentage_asc / percentage_desc
	Page     int      `json:"page"`       // 頁碼，從 1 開始
	PageSize int      `json:"page_size"`  // 每頁筆數
}

// EquipSearchReq 裝備列表查詢請求
type EquipSearchReq struct {
	Date     string   `json:"date"`       // 查詢日期，格式 2006-01-02，空值預設今天
	Category []string `json:"category"`   // 分類篩選，空值為全部
	SortBy   string   `json:"sort_by"`    // 排序方式：price_desc / price_asc / change_desc / change_asc
	Page     int      `json:"page"`       // 頁碼，從 1 開始
	PageSize int      `json:"page_size"`  // 每頁筆數
}

// RecordPriceReq 新增/更新價格請求
type RecordPriceReq struct {
	Price float64 `json:"price" binding:"required,gt=0"` // 價格，必填且必須大於 0
	Date  string  `json:"date"`                          // 記錄日期，格式 2006-01-02，空值預設今天
}
