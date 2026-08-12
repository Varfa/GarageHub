package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type OrderCreatePageData struct {
	Clients []models.Client
	Cars    []models.Car
	Error   string
}
type OrderPageData struct {
	Order  models.Order
	Client models.Client
	Car    models.Car
	Notes  []models.OrderNote
}
type OrdersPageData struct {
	Orders []models.OrderListItem

	Error string
}
type OrderHandler struct {
	service       *service.OrderService
	clientService *service.ClientService
	carService    *service.CarService
	noteService   *service.OrderNoteService
}

func NewOrderHandler(
	service *service.OrderService,
	clientService *service.ClientService,
	carService *service.CarService,
	noteService *service.OrderNoteService,
) *OrderHandler {
	return &OrderHandler{
		service:       service,
		clientService: clientService,
		carService:    carService,
		noteService:   noteService,
	}
}
func (h *OrderHandler) Orders(
	w http.ResponseWriter,
	r *http.Request,
) {
	errorMessage := r.URL.Query().Get("error")
	orders, err := h.service.List(r.Context())
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

	data := OrderPageData{
		Order:  *order,
		Car:    *car,
		Client: *client,
		Notes:  notes,
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
	data := OrderCreatePageData{
		Clients: clients,
		Cars:    cars,
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

	if err := h.service.Create(
		r.Context(),
		order,
	); err != nil {
		http.Error(
			w,
			translate(r, "orders.error.internal"),
			http.StatusBadRequest,
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
