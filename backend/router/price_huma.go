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
		Summary:     "Record item price",
		Description: "Scraper or admin writes the latest auction price for an item.",
		Tags:        []string{"prices"},
	}, h.RecordPrice)

	huma.Register(publicApi, huma.Operation{
		OperationID: "get-latest-price",
		Method:      "GET",
		Path:        "/items/{id}/prices/latest",
		Summary:     "Get latest price for an item",
		Description: "Returns the most recent price record for the specified item. Returns 404 if no price has been recorded yet.",
		Tags:        []string{"prices"},
	}, h.GetLatest)

	huma.Register(publicApi, huma.Operation{
		OperationID: "get-latest-prices-batch",
		Method:      "POST",
		Path:        "/items/prices/latest-batch",
		Summary:     "Batch get latest prices for multiple items",
		Description: "Returns a map of item_id → latest price. Item IDs with no record are omitted.",
		Tags:        []string{"prices"},
	}, h.GetLatestBatch)

	// ── 會員路由（/api/v1/member）────────────────────────────────────
	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-scrolls",
		Method:      "POST",
		Path:        "/scrolls/search",
		Summary:     "Search scroll prices with pagination",
		Description: "Query scroll item prices with filtering by category and percentage, sorted by the specified field. Returns a paginated list of price summaries. Requires member authentication.",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetScrollSummary)

	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-skillbooks",
		Method:      "POST",
		Path:        "/skillbooks/search",
		Summary:     "Search skillbook prices with pagination",
		Description: "Query skill book item prices with filtering by category, sorted by the specified field. Defaults to price_desc sort order. Returns a paginated list of price summaries. Requires member authentication.",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetSkillBookSummary)

	huma.Register(memberApi, huma.Operation{
		OperationID: "member-search-equips",
		Method:      "POST",
		Path:        "/equips/search",
		Summary:     "Search equipment prices with pagination",
		Description: "Query equipment item prices with filtering by category, sorted by the specified field. Defaults to price_desc sort order. Returns a paginated list of price summaries. Requires member authentication.",
		Tags:        []string{"prices", "member"},
		Security:    []map[string][]string{{"memberBearerAuth": {}}},
	}, h.GetEquipSummary)

	// ── 管理員路由（/api/v1/admin）───────────────────────────────────
	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-record-price",
		Method:      "POST",
		Path:        "/items/{id}/prices",
		Summary:     "Admin: record item price (requires price:write permission)",
		Description: "Manually record the latest auction price for an item. Requires the caller to have `price:write` Casbin permission (superadmin is exempt). The source is always recorded as `admin`.",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminRecordPrice)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-get-price-history",
		Method:      "GET",
		Path:        "/items/{id}/prices",
		Summary:     "Admin: get all price records for an item",
		Description: "Returns every price record (one per date) stored for the specified item, ordered by date descending. Useful for reviewing the price time series.",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminGetHistory)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-get-price-histories",
		Method:      "GET",
		Path:        "/items/{id}/histories",
		Summary:     "Admin: get full audit history for an item",
		Description: "Returns the full audit log of every write operation made to an item's price, including hidden entries. Each entry shows who recorded the price, from what source, and whether it has been hidden.",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.AdminGetPriceHistories)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-delete-history",
		Method:      "DELETE",
		Path:        "/histories/{id}",
		Summary:     "Admin: permanently delete a price history entry",
		Description: "Permanently removes a single price history entry by its ID. This action is irreversible. Use the hidden flag instead if you want to soft-delete.",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.DeletePriceHistory)

	huma.Register(adminApi, huma.Operation{
		OperationID: "admin-toggle-history-hidden",
		Method:      "PATCH",
		Path:        "/histories/{id}/hidden",
		Summary:     "Admin: toggle the hidden flag on a price history entry",
		Description: "Sets or clears the `is_hidden` flag on a price history entry. Hidden entries are excluded from public-facing price queries but are retained in the audit log.",
		Tags:        []string{"prices", "admin"},
		Security:    []map[string][]string{{"adminBearerAuth": {}}},
	}, h.TogglePriceHistoryHidden)
}
