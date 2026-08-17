package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/middleware"
	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type OrderCreatePageData struct {
	Clients   []models.Client
	Cars      []models.Car
	Employees []models.EmployeeListItem
	Error     string
}

type OrderPageData struct {
	Order              models.Order
	Client             models.Client
	Car                models.Car
	Notes              []models.OrderNote
	Employees          []models.OrderEmployeeListItem
	AvailableEmployees []models.EmployeeListItem
	Error              string
	IsOwner            bool
}

type OrdersPageData struct {
	Orders []models.OrderListItem
	Error  string
	Search string
}
type OrderHandler struct {
	service         *service.OrderService
	clientService   *service.ClientService
	carService      *service.CarService
	noteService     *service.OrderNoteService
	employeeService *service.OrderEmployeeService
	staffService    *service.EmployeeService
}

func NewOrderHandler(
	service *service.OrderService,
	clientService *service.ClientService,
	carService *service.CarService,
	noteService *service.OrderNoteService,
	employeeService *service.OrderEmployeeService,
	staffService *service.EmployeeService,
) *OrderHandler {
	return &OrderHandler{
		service:         service,
		clientService:   clientService,
		carService:      carService,
		noteService:     noteService,
		employeeService: employeeService,
		staffService:    staffService,
	}
}
func (h *OrderHandler) Orders(
	w http.ResponseWriter,
	r *http.Request,
) {
	errorMessage := r.URL.Query().Get("error")
	search := r.URL.Query().Get("search")

	orders, err := h.service.List(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := OrdersPageData{
		Orders: orders,
		Error:  errorMessage,
		Search: search,
	}

	RenderTemplate(
		w,
		r,
		"orders",
		data,
	)
}
func (h *OrderHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {
	errorMessage := r.URL.Query().Get("error")
	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}
	order, err := h.service.GetByID(
		r.Context(),
		id,
	)

	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(
			err,
			service.ErrOrderInvalidID,
		):
			status = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrOrderNotFound,
		):
			status = http.StatusNotFound
		}

		http.Error(
			w,
			translate(r, "orders.error.internal"),
			status,
		)
		return
	}

	client, err := h.clientService.GetByID(r.Context(), order.ClientID)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	car, err := h.carService.GetByID(r.Context(), order.CarID)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	notes, err := h.noteService.ListByOrderID(
		r.Context(),
		int64(order.ID),
	)
	if err != nil {
		http.Error(w, translate(r, "orders.error.internal"), http.StatusInternalServerError)
		return
	}

	employees, err := h.employeeService.ListActiveByOrderID(
		r.Context(),
		int64(order.ID),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}
	availableEmployees, err := h.staffService.ListActive(
		r.Context(),
		"",
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}
	currentUser, ok := middleware.CurrentUser(r)

	isOwner := false
	if ok && currentUser != nil {
		isOwner = currentUser.IsOwner
	}
	data := OrderPageData{
		Order:              *order,
		Car:                *car,
		Client:             *client,
		Notes:              notes,
		Employees:          employees,
		AvailableEmployees: availableEmployees,
		Error:              errorMessage,
		IsOwner:            isOwner,
	}
	RenderTemplate(w, r, "order_view", data)

}
func (h *OrderHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	clients, err := h.clientService.List(r.Context(), "")
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}
	cars, err := h.carService.List(r.Context(), "")
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}
	employees, err := h.staffService.ListActive(
		r.Context(),
		"",
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}
	data := OrderCreatePageData{
		Clients:   clients,
		Cars:      cars,
		Employees: employees,
	}
	RenderTemplate(
		w,
		r,
		"order_create",
		data,
	)
}
func (h *OrderHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	clientID, err := strconv.Atoi(
		r.FormValue("client_id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.client_required"),
			http.StatusBadRequest,
		)
		return
	}

	carID, err := strconv.Atoi(
		r.FormValue("car_id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.car_required"),
			http.StatusBadRequest,
		)
		return
	}

	order := models.Order{
		ClientID:  clientID,
		CarID:     carID,
		Complaint: r.FormValue("complaint"),
		Diagnosis: r.FormValue("diagnosis"),
		Note:      r.FormValue("note"),
	}

	orderID, err := h.service.Create(
		r.Context(),
		order,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusBadRequest,
		)
		return
	}

	// Получаем employee_id из формыOrderCreatePageData.
	// Поле необязательное, пустая строка — это нормально.
	employeeIDValue := r.FormValue("employee_id")

	if employeeIDValue != "" {
		employeeID, err := strconv.ParseInt(
			employeeIDValue,
			10,
			64,
		)
		if err != nil {
			http.Error(
				w,
				translate(r, "orders.error.internal"),
				http.StatusBadRequest,
			)
			return
		}

		err = h.employeeService.Assign(
			r.Context(),
			orderID,
			employeeID,
		)
		if err != nil {
			http.Error(
				w,
				translate(r, "orders.error.internal"),
				http.StatusInternalServerError,
			)
			return
		}
	}
	http.Redirect(
		w,
		r,
		"/orders",
		http.StatusSeeOther,
	)
}
func (h *OrderHandler) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}
	id, err := strconv.Atoi(
		r.FormValue("id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}
	status := r.FormValue("status")
	err = h.service.UpdateStatus(
		r.Context(),
		id,
		status)
	if err != nil {

		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, service.ErrOrderInvalidID),
			errors.Is(err, service.ErrOrderInvalidStatus):

			statusCode = http.StatusBadRequest

		case errors.Is(err, service.ErrOrderNotFound):
			statusCode = http.StatusNotFound
		}

		http.Error(
			w,
			translate(r, "orders.error.internal"),
			statusCode,
		)

		return
	}
	http.Redirect(w, r, "/orders", http.StatusSeeOther)

}
func (h *OrderHandler) AddNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(
		r.FormValue("order_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, translate(r, "order.invalid_id"),
			http.StatusBadRequest)
		return
	}
	note := models.OrderNote{
		OrderID: orderID,
		Text:    r.FormValue("text"),
	}
	err = h.noteService.Create(r.Context(), note)
	if err != nil {
		status := http.StatusInternalServerError

		switch {

		case errors.Is(
			err,
			service.ErrOrderInvalidID,
		):
			status = http.StatusBadRequest
		case errors.Is(
			err,
			service.ErrOrderNoteTextRequired,
		):
			status = http.StatusBadRequest

		}
		http.Error(w, translate(r, "orders.error.internal"), status)
		return
	}
	http.Redirect(
		w,
		r,
		"/orders/view?id="+strconv.FormatInt(orderID, 10), http.StatusSeeOther,
	)
}
func (h *OrderHandler) Closed(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")

	orders, err := h.service.ListClosed(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := OrdersPageData{
		Orders: orders,
		Search: search,
	}

	RenderTemplate(
		w,
		r,
		"orders_closed",
		data,
	)
}

func (h *OrderHandler) Restore(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.Atoi(
		r.FormValue("id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Restore(
		r.Context(),
		id,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(
			err,
			service.ErrOrderInvalidID,
		):
			statusCode = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrOrderNotFound,
		):
			statusCode = http.StatusNotFound
		}

		http.Error(
			w,
			translate(r, "orders.error.internal"),
			statusCode,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/orders",
		http.StatusSeeOther,
	)
}
func (h *OrderHandler) AssignEmployee(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	orderID, err := strconv.ParseInt(
		r.FormValue("order_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.invalid_id"),
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
			translate(r, "orders.error.internal"),
			http.StatusBadRequest,
		)
		return
	}

	err = h.employeeService.Assign(
		r.Context(),
		orderID,
		employeeID,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := translate(
			r,
			"orders.error.internal",
		)

		switch {
		case errors.Is(
			err,
			service.ErrOrderEmployeeInvalidOrderID,
		):
			statusCode = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrOrderEmployeeInvalidEmployeeID,
		):
			statusCode = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrOrderEmployeeAlreadyAssigned,
		):
			http.Redirect(
				w,
				r,
				"/orders/view?id="+strconv.FormatInt(orderID, 10)+
					"&error=employee_already_assigned",
				http.StatusSeeOther,
			)
			return
		}

		http.Error(
			w,
			message,
			statusCode,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/orders/view?id="+strconv.FormatInt(orderID, 10),
		http.StatusSeeOther,
	)
}
func (h *OrderHandler) UnassignEmployee(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "orders.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	orderID, err := strconv.ParseInt(
		r.FormValue("order_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "order.invalid_id"),
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
			translate(r, "orders.error.internal"),
			http.StatusBadRequest,
		)
		return
	}

	err = h.employeeService.Unassign(
		r.Context(),
		orderID,
		employeeID,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(
			err,
			service.ErrOrderEmployeeInvalidOrderID,
		):
			statusCode = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrOrderEmployeeInvalidEmployeeID,
		):
			statusCode = http.StatusBadRequest
		}

		http.Error(
			w,
			translate(r, "orders.error.internal"),
			statusCode,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/orders/view?id="+strconv.FormatInt(orderID, 10),
		http.StatusSeeOther,
	)
}
func (h *OrderHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}
	idString := r.FormValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Delete(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/orders",
		http.StatusSeeOther,
	)
}
