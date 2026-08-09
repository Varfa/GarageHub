package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type WarehouseArchivePageData struct {
	Items   []models.WarehouseItem
	Search  string
	Error   string
	Success string
}
type WarehouseViewPageData struct {
	Item    *models.WarehouseItem
	Error   string
	Success string
}
type WarehouseHandler struct {
	service *service.WarehouseService
}

type WarehouseCreatePageData struct {
	Error string
}

type WarehousePageData struct {
	Items   []models.WarehouseItem
	Search  string
	Error   string
	Success string
}

func NewWarehouseHandler(
	service *service.WarehouseService,
) *WarehouseHandler {
	return &WarehouseHandler{
		service: service,
	}
}

func (h *WarehouseHandler) Warehouse(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")
	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	items, err := h.service.List(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := WarehousePageData{
		Items:   items,
		Search:  search,
		Error:   errorMessage,
		Success: successMessage,
	}

	RenderTemplate(
		w,
		r,
		"warehouse",
		data,
	)
}

func (h *WarehouseHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	errorMessage := r.URL.Query().Get("error")

	data := WarehouseCreatePageData{
		Error: errorMessage,
	}

	RenderTemplate(
		w,
		r,
		"warehouse_create",
		data,
	)
}

func (h *WarehouseHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
			http.StatusBadRequest,
		)
		return
	}

	quantity, err := strconv.Atoi(
		r.FormValue("quantity"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректное количество",
			http.StatusBadRequest,
		)
		return
	}

	minQuantity, err := strconv.Atoi(
		r.FormValue("min_quantity"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректное минимальное количество",
			http.StatusBadRequest,
		)
		return
	}

	salePriceCents, err := parseMoneyToCents(
		r.FormValue("sale_price"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректная цена продажи",
			http.StatusBadRequest,
		)
		return
	}

	purchasePriceCents, err := parseMoneyToCents(
		r.FormValue("purchase_price"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректная цена закупки",
			http.StatusBadRequest,
		)
		return
	}

	item := models.WarehouseItem{
		Name:               r.FormValue("name"),
		SKU:                r.FormValue("sku"),
		Manufacturer:       r.FormValue("manufacturer"),
		PurchasePriceCents: purchasePriceCents,
		SalePriceCents:     salePriceCents,
		Quantity:           quantity,
		MinQuantity:        minQuantity,
		Location:           r.FormValue("location"),
		Note:               r.FormValue("note"),
	}

	if err := h.service.Create(
		r.Context(),
		item,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/warehouse/create?error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"складская позиция успешно создана",
	)

	http.Redirect(
		w,
		r,
		"/warehouse?success="+successMessage,
		http.StatusSeeOther,
	)
}
func (h *WarehouseHandler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(
		r.URL.Query().Get("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id ",
			http.StatusBadRequest,
		)
		return
	}
	item, err := h.service.GetByID(r.Context(),
		id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNotFound,
		)
		return
	}
	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")
	data := WarehouseViewPageData{
		Item:    item,
		Error:   errorMessage,
		Success: successMessage,
	}

	RenderTemplate(
		w,
		r,
		"warehouse_item",
		data,
	)
}
func (h *WarehouseHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id складской позиции",
			http.StatusBadRequest,
		)
		return
	}

	quantity, err := strconv.Atoi(
		r.FormValue("quantity"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректное количество",
			http.StatusBadRequest,
		)
		return
	}

	minQuantity, err := strconv.Atoi(
		r.FormValue("min_quantity"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректное минимальное количество",
			http.StatusBadRequest,
		)
		return
	}

	salePriceCents, err := parseMoneyToCents(
		r.FormValue("sale_price"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректная цена продажи",
			http.StatusBadRequest,
		)
		return
	}

	purchasePriceCents, err := parseMoneyToCents(
		r.FormValue("purchase_price"),
	)
	if err != nil {
		http.Error(
			w,
			"некорректная цена закупки",
			http.StatusBadRequest,
		)
		return
	}

	item := models.WarehouseItem{
		ID:                 id,
		Name:               r.FormValue("name"),
		SKU:                r.FormValue("sku"),
		Manufacturer:       r.FormValue("manufacturer"),
		PurchasePriceCents: purchasePriceCents,
		SalePriceCents:     salePriceCents,
		Quantity:           quantity,
		MinQuantity:        minQuantity,
		Location:           r.FormValue("location"),
		Note:               r.FormValue("note"),
	}

	if err := h.service.Update(
		r.Context(),
		item,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/warehouse/view?id="+
				strconv.FormatInt(id, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"складская позиция успешно обновлена",
	)

	http.Redirect(
		w,
		r,
		"/warehouse/view?id="+
			strconv.FormatInt(id, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}
func (h *WarehouseHandler) Archive(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id складской позиции",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Archive(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/warehouse/view?id="+
				strconv.FormatInt(id, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"складская позиция перемещена в архив",
	)

	http.Redirect(
		w,
		r,
		"/warehouse?success="+successMessage,
		http.StatusSeeOther,
	)
}

func (h *WarehouseHandler) Restore(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id складской позиции",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Restore(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/warehouse/view?id="+
				strconv.FormatInt(id, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"складская позиция восстановлена",
	)

	http.Redirect(
		w,
		r,
		"/warehouse/view?id="+
			strconv.FormatInt(id, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}
func (h *WarehouseHandler) ArchivePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")
	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	items, err := h.service.ListArchived(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := WarehouseArchivePageData{
		Items:   items,
		Search:  search,
		Error:   errorMessage,
		Success: successMessage,
	}

	RenderTemplate(
		w,
		r,
		"warehouse_archive",
		data,
	)
}
func parseMoneyToCents(value string) (int64, error) {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, ",", ".")

	if value == "" {
		return 0, nil
	}

	parts := strings.Split(value, ".")

	if len(parts) > 2 {
		return 0, errors.New("некорректная цена")
	}

	euros, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, errors.New("некорректная цена")
	}

	var cents int64

	if len(parts) == 2 {
		fraction := parts[1]

		if len(fraction) == 1 {
			fraction += "0"
		}

		if len(fraction) > 2 {
			return 0, errors.New("некорректная цена")
		}

		cents, err = strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return 0, errors.New("некорректная цена")
		}
	}

	return euros*100 + cents, nil
}

func formatMoney(cents int64) string {
	euros := cents / 100
	remainder := cents % 100

	return fmt.Sprintf(
		"%d.%02d",
		euros,
		remainder,
	)
}
