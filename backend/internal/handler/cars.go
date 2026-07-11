package handler

import (
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type CarHandler struct {
	service       *service.CarService
	clientService *service.ClientService
}
type CarsPageData struct {
	Cars   []models.Car
	Search string
}
type CarPageData struct {
	Car    models.Car
	Client models.Client
	Edit   bool
}
type CarCreatePageData struct {
	Clients []models.Client
	Error   string
}

func NewCarHandler(
	service *service.CarService,
	clientService *service.ClientService,
) *CarHandler {
	return &CarHandler{
		service:       service,
		clientService: clientService,
	}
}

func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		http.Error(w, "некорректный id клиента", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "некорректный год автомобиля", http.StatusBadRequest)
		return
	}

	powerKW, err := strconv.Atoi(r.FormValue("power_kw"))
	if err != nil {
		http.Error(w, "некорректная мощность автомобиля", http.StatusBadRequest)
		return
	}

	mileage, err := strconv.Atoi(r.FormValue("mileage"))
	if err != nil {
		http.Error(w, "некорректный пробег автомобиля", http.StatusBadRequest)
		return
	}

	car := models.Car{
		ClientID:    clientID,
		Brand:       r.FormValue("brand"),
		Model:       r.FormValue("model"),
		Year:        year,
		VIN:         r.FormValue("vin"),
		PlateNumber: r.FormValue("plate_number"),
		Engine:      r.FormValue("engine"),
		PowerKW:     powerKW,
		Color:       r.FormValue("color"),
		Mileage:     mileage,
		Note:        r.FormValue("note"),
	}

	if err := h.service.Create(r.Context(), car); err != nil {
		clients, listErr := h.clientService.List(r.Context(), "")
		if listErr != nil {
			http.Error(w, listErr.Error(), http.StatusInternalServerError)
			return
		}

		data := CarCreatePageData{
			Clients: clients,
			Error:   err.Error(),
		}

		RenderTemplate(w, "car_create", data)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients/view?id="+strconv.Itoa(clientID),
		http.StatusSeeOther,
	)
	http.Redirect(w, r, "/clients/view?id="+strconv.Itoa(clientID), http.StatusSeeOther)
}
func (h *CarHandler) Cars(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	cars, err := h.service.List(r.Context(), search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CarsPageData{
		Cars:   cars,
		Search: search,
	}

	RenderTemplate(w, "cars", data)
}

func (h *CarHandler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "некорректный id автомобиля", http.StatusBadRequest)
		return
	}

	car, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	client, err := h.clientService.GetByID(r.Context(), car.ClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := CarPageData{
		Car:    *car,
		Client: *client,
		Edit:   r.URL.Query().Get("edit") == "1",
	}

	RenderTemplate(w, "car", data)
}
func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "некорректный id автомобиля", http.StatusBadRequest)
		return
	}
	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		http.Error(w, "некорректный id клиента", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "некорректный год выпуска автомобиля", http.StatusBadRequest)
		return
	}
	powerKW, err := strconv.Atoi(r.FormValue("power_kw"))
	if err != nil {
		http.Error(w, "некорректная мощность автомобиля", http.StatusBadRequest)
		return
	}
	mileage, err := strconv.Atoi(r.FormValue("mileage"))
	if err != nil {
		http.Error(w, "некорректный пробег автомобиля", http.StatusBadRequest)
		return
	}
	car := models.Car{
		ID:          id,
		ClientID:    clientID,
		Brand:       r.FormValue("brand"),
		Model:       r.FormValue("model"),
		Year:        year,
		VIN:         r.FormValue("vin"),
		PlateNumber: r.FormValue("plate_number"),
		Engine:      r.FormValue("engine"),
		Color:       r.FormValue("color"),
		Note:        r.FormValue("note"),
		PowerKW:     powerKW,
		Mileage:     mileage,
	}

	if err := h.service.Update(r.Context(), car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/cars/view?id="+strconv.Itoa(id), http.StatusSeeOther)
}
func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "некорректный id автомобиля", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/cars", http.StatusSeeOther)
}
func (h *CarHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	clients, err := h.clientService.List(r.Context(), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CarCreatePageData{
		Clients: clients,
	}

	RenderTemplate(w, "car_create", data)
}
