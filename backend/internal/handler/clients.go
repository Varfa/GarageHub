package handler

import (
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
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := ClientsPageData{
		Clients: clients,
		Search:  search,
		Error:   errorMessage,
	}

	RenderTemplate(w, r, "clients", data)
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
			"метод не разрешён",
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

	client := models.Client{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		Address: r.FormValue("address"),
		Note:    r.FormValue("note"),
	}

	if err := h.service.Create(r.Context(), client); err != nil {
		data := ClientCreatePageData{
			Error: err.Error(),
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(
			w,
			"некорректный id клиента",
			http.StatusBadRequest,
		)
		return
	}

	client, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNotFound,
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
			err.Error(),
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
			err.Error(),
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
			"метод не разрешён",
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

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(
			w,
			"некорректный id клиента",
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

	if err := h.service.Update(r.Context(), client); err != nil {
		http.Error(
			w,
			err.Error(),
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
			"метод не разрешён",
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

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(
			w,
			"некорректный id клиента",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		errorMessage := url.QueryEscape(err.Error())

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
