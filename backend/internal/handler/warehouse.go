package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
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
			translate(r, "warehouse.error.internal"),
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
			translate(r, "warehouse.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "warehouse.request.invalid"),
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
			translate(r, "warehouse.validation.invalid_quantity"),
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
			translate(r, "warehouse.validation.invalid_min_quantity"),
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
			translate(r, "warehouse.validation.invalid_sale_price"),
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
			translate(r, "warehouse.validation.invalid_purchase_price"),
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
			warehouseErrorMessage(r, err),
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
		translate(r, "warehouse.success.created"),
	)

	http.Redirect(
		w,
		r,
		"/warehouse?success="+successMessage,
		http.StatusSeeOther,
	)
}

func (h *WarehouseHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.URL.Query().Get("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "warehouse.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	item, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrItemNotFound) {
			status = http.StatusNotFound
		}

		if errors.Is(err, service.ErrWarehouseInvalidID) {
			status = http.StatusBadRequest
		}

		http.Error(
			w,
			warehouseErrorMessage(r, err),
			status,
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
			translate(r, "warehouse.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "warehouse.request.invalid"),
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
			translate(r, "warehouse.invalid_id"),
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
			translate(r, "warehouse.validation.invalid_quantity"),
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
			translate(r, "warehouse.validation.invalid_min_quantity"),
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
			translate(r, "warehouse.validation.invalid_sale_price"),
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
			translate(r, "warehouse.validation.invalid_purchase_price"),
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
			warehouseErrorMessage(r, err),
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
		translate(r, "warehouse.success.updated"),
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
			translate(r, "warehouse.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "warehouse.request.invalid"),
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
			translate(r, "warehouse.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Archive(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			warehouseErrorMessage(r, err),
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
		translate(r, "warehouse.success.archived"),
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
			translate(r, "warehouse.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "warehouse.request.invalid"),
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
			translate(r, "warehouse.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Restore(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			warehouseErrorMessage(r, err),
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
		translate(r, "warehouse.success.restored"),
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
			translate(r, "warehouse.error.internal"),
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

func warehouseErrorMessage(
	r *http.Request,
	err error,
) string {
	switch {
	case errors.Is(
		err,
		service.ErrWarehouseInvalidID,
	):
		return translate(
			r,
			"warehouse.invalid_id",
		)

	case errors.Is(
		err,
		service.ErrWarehouseNameRequired,
	):
		return translate(
			r,
			"warehouse.validation.name_required",
		)

	case errors.Is(
		err,
		service.ErrWarehouseSKURequired,
	):
		return translate(
			r,
			"warehouse.validation.sku_required",
		)

	case errors.Is(
		err,
		service.ErrWarehousePurchasePriceNegative,
	):
		return translate(
			r,
			"warehouse.validation.purchase_price_negative",
		)

	case errors.Is(
		err,
		service.ErrWarehouseSalePriceNegative,
	):
		return translate(
			r,
			"warehouse.validation.sale_price_negative",
		)

	case errors.Is(
		err,
		service.ErrWarehouseQuantityNegative,
	):
		return translate(
			r,
			"warehouse.validation.quantity_negative",
		)

	case errors.Is(
		err,
		service.ErrWarehouseMinQuantityNegative,
	):
		return translate(
			r,
			"warehouse.validation.min_quantity_negative",
		)

	case errors.Is(
		err,
		service.ErrWarehouseAlreadyArchived,
	):
		return translate(
			r,
			"warehouse.already_archived",
		)

	case errors.Is(
		err,
		service.ErrWarehouseAlreadyActive,
	):
		return translate(
			r,
			"warehouse.already_active",
		)

	case errors.Is(
		err,
		repository.ErrItemNotFound,
	):
		return translate(
			r,
			"warehouse.not_found",
		)

	default:
		return translate(
			r,
			"warehouse.error.internal",
		)
	}
}

func parseMoneyToCents(
	value string,
) (int64, error) {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(
		value,
		",",
		".",
	)

	if value == "" {
		return 0, nil
	}

	parts := strings.Split(
		value,
		".",
	)

	if len(parts) > 2 {
		return 0, errors.New(
			"invalid money value",
		)
	}

	euros, err := strconv.ParseInt(
		parts[0],
		10,
		64,
	)
	if err != nil {
		return 0, errors.New(
			"invalid money value",
		)
	}

	var cents int64

	if len(parts) == 2 {
		fraction := parts[1]

		if len(fraction) == 1 {
			fraction += "0"
		}

		if len(fraction) > 2 {
			return 0, errors.New(
				"invalid money value",
			)
		}

		cents, err = strconv.ParseInt(
			fraction,
			10,
			64,
		)
		if err != nil {
			return 0, errors.New(
				"invalid money value",
			)
		}
	}

	return euros*100 + cents, nil
}
