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
	svc         service.ItemService
	queryRepo   repository.QueryRepository
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
	filterType, _ := strconv.Atoi(c.DefaultQuery("filter_type", "0"))
	filterPriority, _ := strconv.Atoi(c.DefaultQuery("filter_priority", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}

	result, err := h.svc.GetAllWithLatestPrice(sortBy, search, filterType, filterPriority, page, pageSize)
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
	itemType, _ := strconv.Atoi(c.DefaultQuery("item_type", "0"))
	categories, err := h.categoryRepo.FindByItemType(itemType)
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
	today := time.Now().Format("2006-01-02")
	go func() { _ = h.queryRepo.RecordItemView(id, today) }()
	respOK(c, summary)
}

func (h *ItemHandler) GetTracked(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
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

func (h *ItemHandler) buildExcel(dates [7]string) (*excelize.File, string, error) {
	equips, err := h.svc.GetAllForExport(int(model.ItemTypeEquip), dates)
	if err != nil {
		return nil, "", err
	}
	skillbooks, err := h.svc.GetAllForExport(int(model.ItemTypeSkillBook), dates)
	if err != nil {
		return nil, "", err
	}
	scrolls, err := h.svc.GetAllForExport(int(model.ItemTypeScroll), dates)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()

	numStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 3})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
		Font: &excelize.Font{Bold: true},
	})

	priceFields := func(row model.ExportRow) [7]*float64 {
		return [7]*float64{row.D0Price, row.D1Price, row.D2Price, row.D3Price, row.D4Price, row.D5Price, row.D6Price}
	}

	sheets := []struct {
		name string
		rows []model.ExportRow
	}{
		{"裝備", equips},
		{"技能書", skillbooks},
		{"卷軸", scrolls},
	}

	headers := []string{"名稱", "類型"}
	for j, d := range dates {
		t, _ := time.Parse("2006-01-02", d)
		label := t.Format("1/2")
		if j == 0 {
			label = fmt.Sprintf("今日(%s)", label)
		}
		headers = append(headers, label)
	}
	lastCol, _ := excelize.ColumnNumberToName(len(headers))

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
			for col, p := range priceFields(row) {
				if p != nil {
					cell := mustCell(col+3, r)
					f.SetCellValue(ws, cell, int64(*p))
					f.SetCellStyle(ws, cell, cell, numStyle)
				}
			}
		}

		f.SetColWidth(ws, "A", "A", 36)
		f.SetColWidth(ws, "B", "B", 14)
		f.SetColWidth(ws, "C", lastCol, 14)
	}

	filename := fmt.Sprintf("artale_market_%s.xlsx", dates[0])
	return f, filename, nil
}

func (h *ItemHandler) ExportExcel(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	ref := time.Now().In(loc)
	var dates [7]string
	for i := range dates {
		dates[i] = ref.AddDate(0, 0, -i).Format("2006-01-02")
	}

	f, filename, err := h.buildExcel(dates)
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

	loc, _ := time.LoadLocation("Asia/Taipei")
	ref := time.Now().In(loc)
	var dates [7]string
	for i := range dates {
		dates[i] = ref.AddDate(0, 0, -i).Format("2006-01-02")
	}

	f, filename, err := h.buildExcel(dates)
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
	w.WriteField("payload_json", fmt.Sprintf(`{"content":"📊 Artale Market 價格報表 %s"}`, dates[0]))
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

