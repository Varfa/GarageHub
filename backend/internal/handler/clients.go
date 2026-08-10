package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type ClientHandler struct {
	service    *service.ClientService
	carService *service.CarService
}

type ClientsPageData struct {
	Clients []models.ClientListItem
	Search  string
	Error   string
}

type ClientCreatePageData struct {
	Error string
}

type ClientPageData struct {
	Client     models.Client
	Cars       []models.Car
	AllClients []models.Client
	Edit       bool
}

func NewClientHandler(
	service *service.ClientService,
	carService *service.CarService,
) *ClientHandler {
	return &ClientHandler{
		service:    service,
		carService: carService,
	}
}

func (h *ClientHandler) Clients(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")
	errorMessage := r.URL.Query().Get("error")

	clients, err := h.service.ListWithCarsCount(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "clients.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := ClientsPageData{
		Clients: clients,
		Search:  search,
		Error:   errorMessage,
	}

	RenderTemplate(
		w,
		r,
		"clients",
		data,
	)
}

func (h *ClientHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	RenderTemplate(
		w,
		r,
		"client_create",
		ClientCreatePageData{},
	)
}

func (h *ClientHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(
				r,
				"clients.request.method_not_allowed",
			),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(
				r,
				"clients.request.invalid",
			),
			http.StatusBadRequest,
		)
		return
	}

	client := models.Client{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		Address: r.FormValue("address"),
		Note:    r.FormValue("note"),
	}

	if err := h.service.Create(
		r.Context(),
		client,
	); err != nil {
		data := ClientCreatePageData{
			Error: clientErrorMessage(
				r,
				err,
			),
		}

		RenderTemplate(
			w,
			r,
			"client_create",
			data,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients",
		http.StatusSeeOther,
	)
}

func (h *ClientHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(
				r,
				"client.invalid_id",
			),
			http.StatusBadRequest,
		)
		return
	}

	client, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(
			err,
			service.ErrClientInvalidID,
		):
			status = http.StatusBadRequest

		case errors.Is(
			err,
			service.ErrClientNotFound,
		):
			status = http.StatusNotFound
		}

		http.Error(
			w,
			clientErrorMessage(
				r,
				err,
			),
			status,
		)
		return
	}

	cars, err := h.carService.ListByClientID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			translate(
				r,
				"clients.error.internal",
			),
			http.StatusInternalServerError,
		)
		return
	}

	allClients, err := h.service.List(
		r.Context(),
		"",
	)
	if err != nil {
		http.Error(
			w,
			translate(
				r,
				"clients.error.internal",
			),
			http.StatusInternalServerError,
		)
		return
	}

	data := ClientPageData{
		Client:     *client,
		Cars:       cars,
		AllClients: allClients,
		Edit:       r.URL.Query().Get("edit") == "1",
	}

	RenderTemplate(
		w,
		r,
		"client",
		data,
	)
}

func (h *ClientHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(
				r,
				"clients.request.method_not_allowed",
			),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(
				r,
				"clients.request.invalid",
			),
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
			translate(
				r,
				"client.invalid_id",
			),
			http.StatusBadRequest,
		)
		return
	}

	client := models.Client{
		ID:      id,
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		Address: r.FormValue("address"),
		Note:    r.FormValue("note"),
	}

	if err := h.service.Update(
		r.Context(),
		client,
	); err != nil {
		http.Error(
			w,
			clientErrorMessage(
				r,
				err,
			),
			http.StatusBadRequest,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients/view?id="+strconv.Itoa(id),
		http.StatusSeeOther,
	)
}

func (h *ClientHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(
				r,
				"clients.request.method_not_allowed",
			),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(
				r,
				"clients.request.invalid",
			),
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
			translate(
				r,
				"client.invalid_id",
			),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Delete(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			clientErrorMessage(
				r,
				err,
			),
		)

		http.Redirect(
			w,
			r,
			"/clients?error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients",
		http.StatusSeeOther,
	)
}

func clientErrorMessage(
	r *http.Request,
	err error,
) string {
	switch {
	case errors.Is(
		err,
		service.ErrClientInvalidID,
	):
		return translate(
			r,
			"client.invalid_id",
		)

	case errors.Is(
		err,
		service.ErrClientNotFound,
	):
		return translate(
			r,
			"client.error.not_found",
		)

	case errors.Is(
		err,
		service.ErrClientHasCars,
	):
		return translate(
			r,
			"client.error.has_cars",
		)

	case errors.Is(
		err,
		service.ErrClientNameRequired,
	):
		return translate(
			r,
			"client.validation.name_required",
		)

	case errors.Is(
		err,
		service.ErrClientPhoneRequired,
	):
		return translate(
			r,
			"client.validation.phone_required",
		)

	case errors.Is(
		err,
		service.ErrClientAlreadyExists,
	):
		return translate(
			r,
			"client.already_exists",
		)

	default:
		return translate(
			r,
			"clients.error.internal",
		)
	}
}
