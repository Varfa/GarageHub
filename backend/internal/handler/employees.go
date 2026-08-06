package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

type EmployeeCreatePageData struct {
	Positions []models.EmployeePosition
	Error     string
}

type EmployeesPageData struct {
	Employees []models.EmployeeListItem
	Search    string
	Error     string
	Success   string
}

type EmployeeArchivePageData struct {
	Employees []models.EmployeeListItem
	Search    string
	Error     string
	Success   string
}

type EmployeePageData struct {
	Employee  models.Employee
	Phones    []models.EmployeePhone
	Positions []models.EmployeePosition
	Edit      bool
	Error     string
	Success   string
}

func NewEmployeeHandler(
	service *service.EmployeeService,
) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) Employees(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")

	employees, err := h.service.ListActive(
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

	data := EmployeesPageData{
		Employees: employees,
		Search:    search,
		Error:     r.URL.Query().Get("error"),
		Success:   r.URL.Query().Get("success"),
	}

	RenderTemplate(
		w,
		r,
		"employees",
		data,
	)
}

func (h *EmployeeHandler) ArchivePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")

	employees, err := h.service.ListArchived(
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

	data := EmployeeArchivePageData{
		Employees: employees,
		Search:    search,
		Error:     r.URL.Query().Get("error"),
		Success:   r.URL.Query().Get("success"),
	}

	RenderTemplate(
		w,
		r,
		"employee_archive",
		data,
	)
}

func (h *EmployeeHandler) View(
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
			"некорректный id сотрудника",
			http.StatusBadRequest,
		)
		return
	}

	employee, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNotFound,
		)
		return
	}

	phones, err := h.service.ListPhones(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	positions, err := h.service.ListPositions(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := EmployeePageData{
		Employee:  *employee,
		Phones:    phones,
		Positions: positions,
		Edit: employee.IsActive &&
			r.URL.Query().Get("edit") == "1",
		Error:   r.URL.Query().Get("error"),
		Success: r.URL.Query().Get("success"),
	}

	RenderTemplate(
		w,
		r,
		"employee",
		data,
	)
}

func (h *EmployeeHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	positions, err := h.service.ListPositions(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := EmployeeCreatePageData{
		Positions: positions,
	}

	RenderTemplate(
		w,
		r,
		"employee_create",
		data,
	)
}

func (h *EmployeeHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
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

	positionID, err := strconv.ParseInt(
		r.FormValue("position_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id должности",
			http.StatusBadRequest,
		)
		return
	}

	employee := models.Employee{
		FirstName:  r.FormValue("first_name"),
		LastName:   r.FormValue("last_name"),
		Phone:      r.FormValue("phone"),
		PositionID: positionID,
	}

	if err := h.service.Create(
		r.Context(),
		employee,
	); err != nil {
		positions, listErr := h.service.ListPositions(
			r.Context(),
		)
		if listErr != nil {
			http.Error(
				w,
				listErr.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		data := EmployeeCreatePageData{
			Positions: positions,
			Error:     err.Error(),
		}

		RenderTemplate(
			w,
			r,
			"employee_create",
			data,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/employees",
		http.StatusSeeOther,
	)
}

func (h *EmployeeHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
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

	employeeID, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id сотрудника",
			http.StatusBadRequest,
		)
		return
	}

	positionID, err := strconv.ParseInt(
		r.FormValue("position_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id должности",
			http.StatusBadRequest,
		)
		return
	}

	currentEmployee, err := h.service.GetByID(
		r.Context(),
		employeeID,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNotFound,
		)
		return
	}

	var email *string

	emailValue := strings.TrimSpace(
		r.FormValue("email"),
	)

	if emailValue != "" {
		email = &emailValue
	}

	employee := models.Employee{
		ID:         employeeID,
		FirstName:  r.FormValue("first_name"),
		LastName:   r.FormValue("last_name"),
		Phone:      currentEmployee.Phone,
		Email:      email,
		PositionID: positionID,
		IsActive:   currentEmployee.IsActive,
		CreatedAt:  currentEmployee.CreatedAt,
		UpdatedAt:  currentEmployee.UpdatedAt,
	}

	if err := h.service.Update(
		r.Context(),
		employee,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/employees/view?id="+
				strconv.FormatInt(employeeID, 10)+
				"&edit=1&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"данные сотрудника успешно обновлены",
	)

	http.Redirect(
		w,
		r,
		"/employees/view?id="+
			strconv.FormatInt(employeeID, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}

func (h *EmployeeHandler) AddPhone(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
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

	employeeID, err := strconv.ParseInt(
		r.FormValue("employee_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id сотрудника",
			http.StatusBadRequest,
		)
		return
	}

	phone := models.EmployeePhone{
		EmployeeID: employeeID,
		Phone:      r.FormValue("phone"),
		Label:      r.FormValue("label"),
		IsPrimary:  r.FormValue("is_primary") == "on",
	}

	if err := h.service.AddPhone(
		r.Context(),
		phone,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/employees/view?id="+
				strconv.FormatInt(employeeID, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"телефон успешно добавлен",
	)

	http.Redirect(
		w,
		r,
		"/employees/view?id="+
			strconv.FormatInt(employeeID, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}

func (h *EmployeeHandler) Archive(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
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

	employeeID, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id сотрудника",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Archive(
		r.Context(),
		employeeID,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/employees/view?id="+
				strconv.FormatInt(employeeID, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"сотрудник перенесён в архив",
	)

	http.Redirect(
		w,
		r,
		"/employees/archive?success="+successMessage,
		http.StatusSeeOther,
	)
}

func (h *EmployeeHandler) Restore(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
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

	employeeID, err := strconv.ParseInt(
		r.FormValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"некорректный id сотрудника",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Restore(
		r.Context(),
		employeeID,
	); err != nil {
		errorMessage := url.QueryEscape(
			err.Error(),
		)

		http.Redirect(
			w,
			r,
			"/employees/view?id="+
				strconv.FormatInt(employeeID, 10)+
				"&error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	successMessage := url.QueryEscape(
		"сотрудник восстановлен",
	)

	http.Redirect(
		w,
		r,
		"/employees?success="+successMessage,
		http.StatusSeeOther,
	)
}
