package handler

import (
	"context"
	"time"

	"artale_market/dto"
	"artale_market/middleware"
	"artale_market/model"
	"artale_market/service"

	"github.com/casbin/casbin/v3"
	"github.com/danielgtaylor/huma/v2"
)

type PriceHumaHandler struct {
	svc      service.PriceService
	querySvc service.QueryService
	enforcer *casbin.Enforcer
}

func NewPriceHumaHandler(svc service.PriceService, querySvc service.QueryService, enforcer *casbin.Enforcer) *PriceHumaHandler {
	return &PriceHumaHandler{svc: svc, querySvc: querySvc, enforcer: enforcer}
}

// ── Input / Output types ─────────────────────────────────────────────────────

type HumaRecordPriceInput struct {
	ItemID  uint   `path:"id" doc:"Item ID"`
	XUserID string `header:"X-User-ID" doc:"Caller identifier for query logging"`
	Body    struct {
		Price float64 `json:"price" doc:"Lowest auction price" minimum:"0"`
		Date  string  `json:"date,omitempty" doc:"Date YYYY-MM-DD, defaults to today (Taipei time)"`
	}
}

type HumaAdminRecordPriceInput struct {
	ItemID uint `path:"id" doc:"Item ID"`
	Body   struct {
		Price float64 `json:"price" doc:"Lowest auction price" minimum:"0"`
		Date  string  `json:"date,omitempty" doc:"Date YYYY-MM-DD, defaults to today (Taipei time)"`
	}
}

type HumaPriceRecordOutput struct {
	Body *model.PriceRecord
}

type HumaGetLatestInput struct {
	ItemID uint `path:"id" doc:"Item ID"`
}

type HumaGetLatestBatchInput struct {
	Body struct {
		ItemIDs []uint `json:"item_ids" doc:"List of item IDs" minItems:"1"`
	}
}

type HumaGetLatestBatchOutput struct {
	Body struct {
		Data map[uint]float64 `json:"data" doc:"Map of item_id to latest price"`
	}
}

type HumaScrollSearchInput struct {
	Body dto.ScrollSearchReq
}

type HumaSkillBookSearchInput struct {
	Body struct {
		Date     string   `json:"date,omitempty"     doc:"查詢日期，格式 YYYY-MM-DD，省略時預設台北今天"                                                                                           format:"date"`
		Category []string `json:"category,omitempty" doc:"職業篩選，省略或空陣列為全部職業。可填：warrior / bowman / magician / thief / pirate 等職業分類名稱"`
		SortBy   string   `json:"sort_by,omitempty"  doc:"排序方式，省略時預設 price_desc。可填：price_desc（價格高→低）/ price_asc（價格低→高）/ change_desc（漲幅高→低）/ change_asc（跌幅高→低）"`
		Page     int      `json:"page,omitempty"     doc:"頁碼，從 1 開始，省略時預設第 1 頁"                                                                                                    minimum:"1"`
		PageSize int      `json:"page_size"          doc:"每頁筆數，必填"                                                                                                                   minimum:"1" maximum:"200"`
	}
}

type HumaEquipSearchInput struct {
	Body struct {
		Date     string   `json:"date,omitempty"     doc:"查詢日期，格式 YYYY-MM-DD，省略時預設台北今天"                                                                                  format:"date"`
		Category []string `json:"category,omitempty" doc:"裝備分類篩選，省略或空陣列為全部。可填：weapon / armor / accessory 等分類名稱"`
		SortBy   string   `json:"sort_by,omitempty"  doc:"排序方式，省略時預設 price_desc。可填：price_desc（價格高→低）/ price_asc（價格低→高）/ change_desc（漲幅高→低）/ change_asc（跌幅高→低）"`
		Page     int      `json:"page,omitempty"     doc:"頁碼，從 1 開始，省略時預設第 1 頁"                                                                                            minimum:"1"`
		PageSize int      `json:"page_size"          doc:"每頁筆數，必填"                                                                                                              minimum:"1" maximum:"200"`
	}
}

type HumaPagedSummaryOutput struct {
	Body *model.PagedSummary
}

// itemID-based inputs (for /items/{id}/... routes)
type HumaItemIDInput struct {
	ItemID uint `path:"id" doc:"Item ID"`
}

// historyID-based inputs (for /histories/{id}/... routes)
type HumaHistoryIDInput struct {
	HistoryID uint `path:"id" doc:"Price history record ID"`
}

type HumaToggleHiddenInput struct {
	HistoryID uint `path:"id" doc:"Price history record ID"`
	Body      struct {
		IsHidden bool `json:"is_hidden" doc:"Whether to hide this history entry"`
	}
}

type HumaPriceRecordListOutput struct {
	Body []model.PriceRecord
}

type HumaPriceHistoryListOutput struct {
	Body []model.PriceHistory
}

type HumaMessageOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// ── Public handlers ──────────────────────────────────────────────────────────

func (h *PriceHumaHandler) RecordPrice(ctx context.Context, input *HumaRecordPriceInput) (*HumaPriceRecordOutput, error) {
	date := input.Body.Date
	if date == "" {
		date = twToday()
	} else if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, huma.Error400BadRequest("date must be YYYY-MM-DD", err)
	}

	// Detect source: admin if caller passes a JWT Authorization header
	source := "scraper"
	if _, ok := ctx.Value(middleware.CtxAdminID{}).(string); ok {
		source = "admin"
	}

	record, err := h.svc.Record(input.ItemID, input.Body.Price, date, source)
	if err != nil {
		return nil, huma.Error404NotFound("item not found", err)
	}

	if input.XUserID != "" {
		go func() { _ = h.querySvc.RecordQuery(input.XUserID, input.ItemID) }()
	}

	return &HumaPriceRecordOutput{Body: record}, nil
}

func (h *PriceHumaHandler) GetLatest(ctx context.Context, input *HumaGetLatestInput) (*HumaPriceRecordOutput, error) {
	record, err := h.svc.GetLatest(input.ItemID)
	if err != nil {
		return nil, huma.Error404NotFound("no price record found", err)
	}
	return &HumaPriceRecordOutput{Body: record}, nil
}

func (h *PriceHumaHandler) GetLatestBatch(ctx context.Context, input *HumaGetLatestBatchInput) (*HumaGetLatestBatchOutput, error) {
	records, err := h.svc.GetLatestBatch(input.Body.ItemIDs)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to query prices", err)
	}
	data := make(map[uint]float64, len(records))
	for _, r := range records {
		data[r.ItemID] = r.Price
	}
	out := &HumaGetLatestBatchOutput{}
	out.Body.Data = data
	return out, nil
}

// ── Member handlers ──────────────────────────────────────────────────────────

func (h *PriceHumaHandler) GetScrollSummary(ctx context.Context, input *HumaScrollSearchInput) (*HumaPagedSummaryOutput, error) {
	b := input.Body
	if b.Date == "" {
		b.Date = twToday()
	}
	if b.Page < 1 {
		b.Page = 1
	}
	if len(b.Category) == 1 && b.Category[0] == "scroll_all" {
		b.Category = nil
	}
	result, err := h.svc.GetScrollSummary(b.Date, b.Percentage, b.Category, b.SortBy, b.Page, b.PageSize)
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed", err)
	}
	return &HumaPagedSummaryOutput{Body: result}, nil
}

func (h *PriceHumaHandler) GetSkillBookSummary(ctx context.Context, input *HumaSkillBookSearchInput) (*HumaPagedSummaryOutput, error) {
	b := input.Body
	if b.Date == "" {
		b.Date = twToday()
	}
	if b.Page < 1 {
		b.Page = 1
	}
	if b.SortBy == "" {
		b.SortBy = "price_desc"
	}
	result, err := h.svc.GetSkillBookSummary(b.Date, b.Category, b.SortBy, b.Page, b.PageSize)
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed", err)
	}
	return &HumaPagedSummaryOutput{Body: result}, nil
}

func (h *PriceHumaHandler) GetEquipSummary(ctx context.Context, input *HumaEquipSearchInput) (*HumaPagedSummaryOutput, error) {
	b := input.Body
	if b.Date == "" {
		b.Date = twToday()
	}
	if b.Page < 1 {
		b.Page = 1
	}
	if b.SortBy == "" {
		b.SortBy = "price_desc"
	}
	result, err := h.svc.GetEquipSummary(b.Date, b.Category, b.SortBy, b.Page, b.PageSize)
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed", err)
	}
	return &HumaPagedSummaryOutput{Body: result}, nil
}

// ── Admin handlers ───────────────────────────────────────────────────────────

func (h *PriceHumaHandler) AdminRecordPrice(ctx context.Context, input *HumaAdminRecordPriceInput) (*HumaPriceRecordOutput, error) {
	// Check price:write permission via Casbin (superadmin bypasses)
	role, _ := ctx.Value(middleware.CtxAdminRole{}).(string)
	if role != "superadmin" {
		username, _ := ctx.Value(middleware.CtxAdminUsername{}).(string)
		ok, err := h.enforcer.Enforce(username, "price", "write")
		if err != nil || !ok {
			return nil, huma.Error403Forbidden("無此操作權限", nil)
		}
	}

	date := input.Body.Date
	if date == "" {
		date = twToday()
	} else if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, huma.Error400BadRequest("date must be YYYY-MM-DD", err)
	}

	record, err := h.svc.Record(input.ItemID, input.Body.Price, date, "admin")
	if err != nil {
		return nil, huma.Error404NotFound("item not found", err)
	}
	return &HumaPriceRecordOutput{Body: record}, nil
}

func (h *PriceHumaHandler) AdminGetHistory(ctx context.Context, input *HumaItemIDInput) (*HumaPriceRecordListOutput, error) {
	records, err := h.svc.GetAllHistory(input.ItemID)
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed", err)
	}
	return &HumaPriceRecordListOutput{Body: records}, nil
}

func (h *PriceHumaHandler) AdminGetPriceHistories(ctx context.Context, input *HumaItemIDInput) (*HumaPriceHistoryListOutput, error) {
	records, err := h.svc.GetPriceHistories(input.ItemID)
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed", err)
	}
	return &HumaPriceHistoryListOutput{Body: records}, nil
}

func (h *PriceHumaHandler) DeletePriceHistory(ctx context.Context, input *HumaHistoryIDInput) (*HumaMessageOutput, error) {
	if err := h.svc.DeletePriceHistory(input.HistoryID); err != nil {
		return nil, huma.Error500InternalServerError("delete failed", err)
	}
	out := &HumaMessageOutput{}
	out.Body.Message = "deleted"
	return out, nil
}

func (h *PriceHumaHandler) TogglePriceHistoryHidden(ctx context.Context, input *HumaToggleHiddenInput) (*HumaMessageOutput, error) {
	if err := h.svc.TogglePriceHistoryHidden(input.HistoryID, input.Body.IsHidden); err != nil {
		return nil, huma.Error500InternalServerError("update failed", err)
	}
	out := &HumaMessageOutput{}
	out.Body.Message = "updated"
	return out, nil
}
