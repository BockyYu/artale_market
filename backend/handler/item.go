package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"artale_market/dto"
	"artale_market/model"
	"artale_market/repository"
	"artale_market/service"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ItemHandler struct {
	svc          service.ItemService
	queryRepo    repository.QueryRepository
	categoryRepo repository.CategoryRepository
}

func NewItemHandler(svc service.ItemService, qr repository.QueryRepository, cr repository.CategoryRepository) *ItemHandler {
	return &ItemHandler{svc: svc, queryRepo: qr, categoryRepo: cr}
}

func (h *ItemHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll()
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"data": items})
}

func (h *ItemHandler) AdminGetAll(c *gin.Context) {
	sortBy := c.Query("sort_by")
	search := c.Query("search")
	filterPriority, _ := strconv.Atoi(c.DefaultQuery("filter_priority", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}

	var filterTypes []int
	for _, s := range c.QueryArray("filter_type") {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			filterTypes = append(filterTypes, v)
		}
	}
	filterCategories := c.QueryArray("filter_category")

	result, err := h.svc.GetAllWithLatestPrice(sortBy, search, filterTypes, filterCategories, filterPriority, page, pageSize)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, result)
}

func (h *ItemHandler) Create(c *gin.Context) {
	var item model.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.Create(&item); err != nil {
		if err == service.ErrItemAlreadyExists {
			respBadRequest(c, err)
		} else {
			respInternal(c, err)
		}
		return
	}
	respCreated(c, item)
}

func (h *ItemHandler) Update(c *gin.Context) {
	var input dto.UpdateItemReq
	if err := c.ShouldBindJSON(&input); err != nil {
		respBadRequest(c, err)
		return
	}
	item, err := h.svc.Update(parseID(c), input.Name, input.EnglishName, input.SearchMode, input.ItemType, input.Percentage, input.Category, input.Description)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, item)
}

func (h *ItemHandler) GetCategories(c *gin.Context) {
	t, _ := strconv.Atoi(c.DefaultQuery("item_type", "0"))
	categories, err := h.categoryRepo.FindByItemType(model.ItemType(t))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, categories)
}

func (h *ItemHandler) GetUsedCategories(c *gin.Context) {
	t, _ := strconv.Atoi(c.DefaultQuery("item_type", "0"))
	categories, err := h.categoryRepo.FindUsedByItemType(model.ItemType(t))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, categories)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(parseID(c)); err != nil {
		respNotFound(c, err)
		return
	}
	respDeleted(c)
}

func (h *ItemHandler) GetByID(c *gin.Context) {
	id := parseID(c)
	summary, err := h.svc.GetPriceSummary(id)
	if err != nil {
		respNotFound(c, err)
		return
	}
	today := twToday()
	go func() { _ = h.queryRepo.RecordItemView(id, today) }()
	respOK(c, summary)
}

func (h *ItemHandler) GetTracked(c *gin.Context) {
	today := twToday()
	items, err := h.svc.GetTracked(today)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, items)
}

func (h *ItemHandler) SetHidden(c *gin.Context) {
	var body struct {
		IsHidden bool `json:"is_hidden"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.SetHidden(parseID(c), body.IsHidden); err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"is_hidden": body.IsHidden})
}

func (h *ItemHandler) SetTracked(c *gin.Context) {
	var input dto.SetTrackedReq
	if err := c.ShouldBindJSON(&input); err != nil {
		respBadRequest(c, err)
		return
	}
	item, err := h.svc.SetTracked(parseID(c), input.TrackPriority)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, item)
}

// hasDataOnDate 回傳指定日期在任意 sheet 中是否有至少一筆價格。
func hasDataOnDate(date string, sheets ...[]model.ExportRowDynamic) bool {
	for _, rows := range sheets {
		for _, row := range rows {
			if row.Prices[date] != nil {
				return true
			}
		}
	}
	return false
}

func (h *ItemHandler) buildExcel() (*excelize.File, string, error) {
	loc := time.FixedZone("Asia/Taipei", 8*60*60)
	today := time.Now().In(loc)
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, loc)

	// 產生本月所有日期（由新到舊）
	var allDates []string
	for d := today; !d.Before(monthStart); d = d.AddDate(0, 0, -1) {
		allDates = append(allDates, d.Format("2006-01-02"))
	}

	equips, err := h.svc.GetAllForExportDynamic(model.ItemTypeEquip, allDates)
	if err != nil {
		return nil, "", err
	}
	skillbooks, err := h.svc.GetAllForExportDynamic(model.ItemTypeSkillBook, allDates)
	if err != nil {
		return nil, "", err
	}
	scrolls, err := h.svc.GetAllForExportDynamic(model.ItemTypeScroll, allDates)
	if err != nil {
		return nil, "", err
	}

	// 過濾掉所有道具都沒有資料的日期
	var dates []string
	for _, d := range allDates {
		if hasDataOnDate(d, equips, skillbooks, scrolls) {
			dates = append(dates, d)
		}
	}

	f := excelize.NewFile()

	numStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 3})
	pctStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 9}) // "0%"，四捨五入到整數
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
		Font: &excelize.Font{Bold: true},
	})

	grayFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
	})
	greenFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"E2EFDA"}},
	})
	blueFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"BDD7EE"}},
	})
	yellowFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFEB9C"}},
	})
	// 漲幅%：上漲顯示紅字，下跌顯示綠字
	riseFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Font: &excelize.Font{Color: "CC0000"},
	})
	fallFmt, _ := f.NewConditionalStyle(&excelize.Style{
		Font: &excelize.Font{Color: "257A34"},
	})

	sheets := []struct {
		name string
		rows []model.ExportRowDynamic
	}{
		{"裝備", equips},
		{"技能書", skillbooks},
		{"卷軸", scrolls},
	}

	// 欄位：名稱、類型、月漲幅%、月差價、[日期、漲幅%] 交錯，最後一天無漲幅欄
	// 月統計固定在 col 3(C) 和 col 4(D)；日期欄從 col 5(E) 開始
	headers := []string{"名稱", "類型", "月漲幅%", "月差價"}
	for j, d := range dates {
		t, _ := time.Parse("2006-01-02", d)
		label := t.Format("1/2")
		if j == 0 {
			label = fmt.Sprintf("今日(%s)", label)
		}
		headers = append(headers, label)
		if j < len(dates)-1 {
			headers = append(headers, "漲幅%")
		}
	}

	for i, sheet := range sheets {
		if i == 0 {
			f.SetSheetName("Sheet1", sheet.name)
		} else {
			f.NewSheet(sheet.name)
		}
		ws := sheet.name

		for col, hdr := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(ws, cell, hdr)
			f.SetCellStyle(ws, cell, cell, headerStyle)
		}

		for rowIdx, row := range sheet.rows {
			r := rowIdx + 2
			f.SetCellValue(ws, mustCell(1, r), row.ItemName)
			f.SetCellValue(ws, mustCell(2, r), row.Category)

			// 月統計：今日價格 vs 本月最早有資料的那天
			todayPrice := row.Prices[dates[0]]
			var monthStartPrice *float64
			for k := len(dates) - 1; k >= 0; k-- {
				if p := row.Prices[dates[k]]; p != nil {
					monthStartPrice = p
					break
				}
			}
			if todayPrice != nil && monthStartPrice != nil && *monthStartPrice > 0 {
				monthPct := (*todayPrice - *monthStartPrice) / *monthStartPrice
				monthDiff := int64(*todayPrice) - int64(*monthStartPrice)
				f.SetCellValue(ws, mustCell(3, r), monthPct)
				f.SetCellStyle(ws, mustCell(3, r), mustCell(3, r), pctStyle)
				f.SetCellValue(ws, mustCell(4, r), monthDiff)
				f.SetCellStyle(ws, mustCell(4, r), mustCell(4, r), numStyle)
			}

			// 每日價格（日期欄從 col 5 開始，+2 偏移）
			for j, date := range dates {
				priceCol := j*2 + 5
				p := row.Prices[date]
				if p != nil {
					cell := mustCell(priceCol, r)
					f.SetCellValue(ws, cell, int64(*p))
					f.SetCellStyle(ws, cell, cell, numStyle)
				}
				if j < len(dates)-1 {
					pctCol := j*2 + 6
					pPrev := row.Prices[dates[j+1]]
					if p != nil && pPrev != nil && *pPrev > 0 {
						change := (*p - *pPrev) / *pPrev
						pctCell := mustCell(pctCol, r)
						f.SetCellValue(ws, pctCell, change)
						f.SetCellStyle(ws, pctCell, pctCell, pctStyle)
					}
				}
			}
		}

		f.SetColWidth(ws, "A", "A", 36)
		f.SetColWidth(ws, "B", "B", 14)
		f.SetColWidth(ws, "C", "C", 9)  // 月漲幅%
		f.SetColWidth(ws, "D", "D", 12) // 月差價
		for j := range dates {
			priceColName, _ := excelize.ColumnNumberToName(j*2 + 5)
			f.SetColWidth(ws, priceColName, priceColName, 12)
			if j < len(dates)-1 {
				pctColName, _ := excelize.ColumnNumberToName(j*2 + 6)
				f.SetColWidth(ws, pctColName, pctColName, 7)
			}
		}

		if len(sheet.rows) > 0 {
			lastRow := len(sheet.rows) + 1
			// 月統計欄條件格式
			for _, col := range []string{"C", "D"} {
				f.SetConditionalFormat(ws, fmt.Sprintf("%s2:%s%d", col, col, lastRow), []excelize.ConditionalFormatOptions{
					{Type: "cell", Criteria: ">", Value: "0", Format: &riseFmt},
					{Type: "cell", Criteria: "<", Value: "0", Format: &fallFmt},
				})
			}
			// 每日欄條件格式
			for j := range dates {
				priceColName, _ := excelize.ColumnNumberToName(j*2 + 5)
				priceRange := fmt.Sprintf("%s2:%s%d", priceColName, priceColName, lastRow)
				f.SetConditionalFormat(ws, priceRange, []excelize.ConditionalFormatOptions{
					{Type: "cell", Criteria: "<", Value: "1000000", Format: &grayFmt},
					{Type: "cell", Criteria: "between", MinValue: "1000000", MaxValue: "9999999", Format: &greenFmt},
					{Type: "cell", Criteria: "between", MinValue: "10000000", MaxValue: "99999999", Format: &blueFmt},
					{Type: "cell", Criteria: ">=", Value: "100000000", Format: &yellowFmt},
				})
				if j < len(dates)-1 {
					pctColName, _ := excelize.ColumnNumberToName(j*2 + 6)
					pctRange := fmt.Sprintf("%s2:%s%d", pctColName, pctColName, lastRow)
					f.SetConditionalFormat(ws, pctRange, []excelize.ConditionalFormatOptions{
						{Type: "cell", Criteria: ">", Value: "0", Format: &riseFmt},
						{Type: "cell", Criteria: "<", Value: "0", Format: &fallFmt},
					})
				}
			}
		}
	}

	filename := fmt.Sprintf("artale_market_%s.xlsx", today.Format("2006-01"))
	return f, filename, nil
}

func (h *ItemHandler) ExportExcel(c *gin.Context) {
	f, filename, err := h.buildExcel()
	if err != nil {
		respInternal(c, err)
		return
	}
	defer f.Close()

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := f.Write(c.Writer); err != nil {
		respInternal(c, err)
	}
}

func (h *ItemHandler) SendExcelToDiscord(c *gin.Context) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		respBadRequest(c, fmt.Errorf("DISCORD_WEBHOOK_URL 未設定"))
		return
	}

	f, filename, err := h.buildExcel()
	if err != nil {
		respInternal(c, err)
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		respInternal(c, err)
		return
	}

	loc := time.FixedZone("Asia/Taipei", 8*60*60)
	month := time.Now().In(loc).Format("2006-01")

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("files[0]", filename)
	if err != nil {
		respInternal(c, err)
		return
	}
	if _, err := io.Copy(fw, &buf); err != nil {
		respInternal(c, err)
		return
	}
	w.WriteField("payload_json", fmt.Sprintf(`{"content":"📊 Artale Market 拍賣價格報表(僅參考用途) %s"}`, month))
	w.Close()

	resp, err := http.Post(webhookURL, w.FormDataContentType(), &body)
	if err != nil {
		respInternal(c, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respInternal(c, fmt.Errorf("Discord 回傳 status %d", resp.StatusCode))
		return
	}

	respOK(c, gin.H{"message": "sent"})
}

func mustCell(col, row int) string {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}
