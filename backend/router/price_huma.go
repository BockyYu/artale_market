package router

import "github.com/danielgtaylor/huma/v2"

// registerPriceHuma 以 Huma v2 格式註冊所有 price handler 路由。
// publicApi  → /api（無需認證）
// memberApi  → /api/v1/member（MemberJWT 保護）
// adminApi   → /api/v1/admin（JWTAuth 保護）
func registerPriceHuma(publicApi, memberApi, adminApi huma.API, d *Deps) {
	h := d.PriceHuma

	// ── 公開路由（/api）──────────────────────────────────────────────
	huma.Register(publicApi, huma.Operation{
		OperationID: "record-price",
		Method:      "POST",
		Path:        "/items/{id}/prices",
		Summary:     "寫入道具價格",
		Description: "由爬蟲或管理員寫入指定道具的最新拍賣價格。若未傳入日期則預設為台北今天。",
		Tags:        []string{"prices"},
	}, h.RecordPrice)

	huma.Register(publicApi, huma.Operation{
		OperationID: "get-latest-price",
		Method:      "GET",
		Path:        "/items/{id}/prices/latest",
		Summary:     "取得道具最新價格",
		Description: "回傳指定道具最新一筆價格紀錄。若該道具尚無任何價格資料則回傳 404。",
		Tags:        []string{"prices"},
	}, h.GetLatest)

	huma.Register(publicApi, huma.Operation{
		OperationID: "get-latest-prices-batch",
		Method:      "POST",
		Path:        "/items/prices/latest-batch",
		Summary:     "批次取得多個道具的最新價格",
		Description: "傳入道具 ID 陣列，回傳 item_id → 最新價格 的對應表。沒有價格紀錄的道具不會出現在結果中。",
		Tags:        []string{"prices"},
	}, h.GetLatestBatch)

	// ── 會員路由（/api/v1/member）────────────────────────────────────
	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-scrolls",
		Method:      "POST",
		Path:        "/v1/member/scrolls/search",
		Summary:     "查詢卷軸價格（分頁）",
		Description: "依分類與成功率篩選卷軸價格，支援多種排序方式，回傳分頁摘要列表。需要會員 JWT 認證。",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetScrollSummary)

	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-skillbooks",
		Method:      "POST",
		Path:        "/v1/member/skillbooks/search",
		Summary:     "查詢技能書價格（分頁）",
		Description: "依職業分類篩選技能書價格，支援多種排序方式，預設以價格高→低排序，回傳分頁摘要列表。需要會員 JWT 認證。",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetSkillBookSummary)

	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-equips",
		Method:      "POST",
		Path:        "/v1/member/equips/search",
		Summary:     "查詢裝備價格（分頁）",
		Description: "依裝備分類篩選裝備價格，支援多種排序方式，預設以價格高→低排序，回傳分頁摘要列表。需要會員 JWT 認證。",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetEquipSummary)

	// ── 管理員路由（/api/v1/admin）───────────────────────────────────
	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-record-price",
		Method:      "POST",
		Path:        "/v1/admin/items/{id}/prices",
		Summary:     "（管理員）手動寫入道具價格",
		Description: "手動寫入指定道具的拍賣價格，來源固定記錄為 admin。需要 Casbin price:write 權限（superadmin 免檢查）。",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminRecordPrice)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-get-price-history",
		Method:      "GET",
		Path:        "/v1/admin/items/{id}/prices",
		Summary:     "（管理員）取得道具所有價格紀錄",
		Description: "回傳指定道具全部的每日價格紀錄，依日期由新到舊排序，可用於檢視價格走勢。",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminGetHistory)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-get-price-histories",
		Method:      "GET",
		Path:        "/v1/admin/items/{id}/histories",
		Summary:     "（管理員）取得道具完整異動記錄",
		Description: "回傳指定道具所有價格寫入操作的完整稽核記錄，包含已隱藏的項目。每筆記錄包含寫入者、來源、是否隱藏等資訊。",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminGetPriceHistories)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-delete-history",
		Method:      "DELETE",
		Path:        "/v1/admin/histories/{id}",
		Summary:     "（管理員）永久刪除價格異動記錄",
		Description: "永久刪除指定 ID 的價格異動記錄，此操作不可逆。若只是不想顯示，建議改用隱藏旗標（PATCH hidden）。",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.DeletePriceHistory)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-toggle-history-hidden",
		Method:      "PATCH",
		Path:        "/v1/admin/histories/{id}/hidden",
		Summary:     "（管理員）切換價格異動記錄的隱藏狀態",
		Description: "設定或取消指定價格異動記錄的 is_hidden 旗標。隱藏的記錄不會出現在前台查詢結果，但仍保留於稽核記錄中。",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.TogglePriceHistoryHidden)
}
