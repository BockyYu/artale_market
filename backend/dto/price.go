package dto

// ScrollSearchReq 卷軸列表查詢請求
type ScrollSearchReq struct {
	Date       string   `json:"date,omitempty"      doc:"查詢日期，格式 YYYY-MM-DD，省略時預設台北今天"                                                                                                  format:"date"`
	Percentage []int    `json:"percentage,omitempty" doc:"成功率篩選，可填 10/20/30/40/50/60/70/80/90/100，空陣列為全部"`
	Category   []string `json:"category,omitempty"  doc:"卷軸分類篩選，傳入 scroll_all 或省略視為全部卷軸。可填：scroll_all / atk_scroll / def_scroll / hp_scroll / 其他分類名稱"`
	SortBy     string   `json:"sort_by,omitempty"   doc:"排序方式，省略時使用預設。可填：price_desc（價格高→低）/ price_asc（價格低→高）/ change_desc（漲幅高→低）/ change_asc（跌幅高→低）/ percentage_asc（成功率低→高）/ percentage_desc（成功率高→低）"`
	Page       int      `json:"page,omitempty"      doc:"頁碼，從 1 開始，省略時預設第 1 頁"                                                                                                          minimum:"1"`
	PageSize   int      `json:"page_size"           doc:"每頁筆數，必填"                                                                                                                        minimum:"1" maximum:"200"`
}

// SkillBookSearchReq 技能書列表查詢請求
type SkillBookSearchReq struct {
	Date     string   `json:"date,omitempty"      doc:"查詢日期，格式 YYYY-MM-DD，省略時預設台北今天"                                                                                           format:"date"`
	Category []string `json:"category,omitempty"  doc:"職業篩選，省略或空陣列為全部職業。可填：warrior / bowman / magician / thief / pirate 等職業分類名稱"`
	SortBy   string   `json:"sort_by,omitempty"   doc:"排序方式，省略時預設 price_desc。可填：price_desc（價格高→低）/ price_asc（價格低→高）/ change_desc（漲幅高→低）/ change_asc（跌幅高→低）"`
	Page     int      `json:"page,omitempty"      doc:"頁碼，從 1 開始，省略時預設第 1 頁"                                                                                                    minimum:"1"`
	PageSize int      `json:"page_size"           doc:"每頁筆數，必填"                                                                                                                   minimum:"1" maximum:"200"`
}

// EquipSearchReq 裝備列表查詢請求
type EquipSearchReq struct {
	Date     string   `json:"date,omitempty"      doc:"查詢日期，格式 YYYY-MM-DD，省略時預設台北今天"                                                                                  format:"date"`
	Category []string `json:"category,omitempty"  doc:"裝備分類篩選，省略或空陣列為全部。可填：weapon / armor / accessory 等分類名稱"`
	SortBy   string   `json:"sort_by,omitempty"   doc:"排序方式，省略時預設 price_desc。可填：price_desc（價格高→低）/ price_asc（價格低→高）/ change_desc（漲幅高→低）/ change_asc（跌幅高→低）"`
	Page     int      `json:"page,omitempty"      doc:"頁碼，從 1 開始，省略時預設第 1 頁"                                                                                            minimum:"1"`
	PageSize int      `json:"page_size"           doc:"每頁筆數，必填"                                                                                                              minimum:"1" maximum:"200"`
}

// RecordPriceReq 新增/更新價格請求
type RecordPriceReq struct {
	Price float64 `json:"price" binding:"required,gt=0"` // 價格，必填且必須大於 0
	Date  string  `json:"date"`                          // 記錄日期，格式 2006-01-02，空值預設今天
}
