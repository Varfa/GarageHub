package handler

import (
	"errors"
	"net/http"
	"net/url"
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
	Error  string
}

type CarPageData struct {
	Car        models.Car
	Client     models.Client
	AllClients []models.Client
	Edit       bool
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

func (h *CarHandler) Cars(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("search")
	errorMessage := r.URL.Query().Get("error")

	cars, err := h.service.List(
		r.Context(),
		search,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "cars.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := CarsPageData{
		Cars:   cars,
		Search: search,
		Error:  errorMessage,
	}

	RenderTemplate(
		w,
		r,
		"cars",
		data,
	)
}

func (h *CarHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	clients, err := h.clientService.List(
		r.Context(),
		"",
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "cars.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := CarCreatePageData{
		Clients: clients,
	}

	RenderTemplate(
		w,
		r,
		"car_create",
		data,
	)
}

func (h *CarHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "cars.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "cars.request.invalid"),
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
			translate(r, "car.client_invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	year, err := strconv.Atoi(
		r.FormValue("year"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.year_invalid"),
			http.StatusBadRequest,
		)
		return
	}

	powerKW, err := strconv.Atoi(
		r.FormValue("power_kw"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.power_invalid"),
			http.StatusBadRequest,
		)
		return
	}

	mileage, err := strconv.Atoi(
		r.FormValue("mileage"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.mileage_invalid"),
			http.StatusBadRequest,
		)
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

	if err := h.service.Create(
		r.Context(),
		car,
	); err != nil {
		clients, listErr := h.clientService.List(
			r.Context(),
			"",
		)
		if listErr != nil {
			http.Error(
				w,
				translate(r, "cars.error.internal"),
				http.StatusInternalServerError,
			)
			return
		}

		data := CarCreatePageData{
			Clients: clients,
			Error:   carErrorMessage(r, err),
		}

		RenderTemplate(
			w,
			r,
			"car_create",
			data,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients/view?id="+strconv.Itoa(clientID),
		http.StatusSeeOther,
	)
}

func (h *CarHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	car, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			carErrorMessage(r, err),
			http.StatusNotFound,
		)
		return
	}

	client, err := h.clientService.GetByID(
		r.Context(),
		car.ClientID,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "cars.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	allClients, err := h.clientService.List(
		r.Context(),
		"",
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "cars.error.internal"),
			http.StatusInternalServerError,
		)
		return
	}

	data := CarPageData{
		Car:        *car,
		Client:     *client,
		AllClients: allClients,
		Edit:       r.URL.Query().Get("edit") == "1",
	}

	RenderTemplate(
		w,
		r,
		"car",
		data,
	)
}

func (h *CarHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "cars.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "cars.request.invalid"),
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
			translate(r, "car.invalid_id"),
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
			translate(r, "car.client_invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	year, err := strconv.Atoi(
		r.FormValue("year"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.year_invalid"),
			http.StatusBadRequest,
		)
		return
	}

	powerKW, err := strconv.Atoi(
		r.FormValue("power_kw"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.power_invalid"),
			http.StatusBadRequest,
		)
		return
	}

	mileage, err := strconv.Atoi(
		r.FormValue("mileage"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.validation.mileage_invalid"),
			http.StatusBadRequest,
		)
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
		PowerKW:     powerKW,
		Color:       r.FormValue("color"),
		Mileage:     mileage,
		Note:        r.FormValue("note"),
	}

	if err := h.service.Update(
		r.Context(),
		car,
	); err != nil {
		http.Error(
			w,
			carErrorMessage(r, err),
			http.StatusBadRequest,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/cars/view?id="+strconv.Itoa(id),
		http.StatusSeeOther,
	)
}

func (h *CarHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "cars.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "cars.request.invalid"),
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
			translate(r, "car.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Delete(
		r.Context(),
		id,
	); err != nil {
		errorMessage := url.QueryEscape(
			carErrorMessage(r, err),
		)

		http.Redirect(
			w,
			r,
			"/cars?error="+errorMessage,
			http.StatusSeeOther,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/cars",
		http.StatusSeeOther,
	)
}

func (h *CarHandler) ChangeOwner(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "cars.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "cars.request.invalid"),
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
			translate(r, "car.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	newClientID, err := strconv.Atoi(
		r.FormValue("client_id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.owner_invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	oldClientID, err := strconv.Atoi(
		r.FormValue("old_client_id"),
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "car.current_owner_invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	returnTo := r.FormValue("return_to")

	if err := h.service.ChangeOwner(
		r.Context(),
		carID,
		newClientID,
	); err != nil {
		http.Error(
			w,
			carErrorMessage(r, err),
			http.StatusBadRequest,
		)
		return
	}

	if returnTo == "car" {
		http.Redirect(
			w,
			r,
			"/cars/view?id="+strconv.Itoa(carID),
			http.StatusSeeOther,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/clients/view?id="+strconv.Itoa(oldClientID),
		http.StatusSeeOther,
	)
}

func carErrorMessage(
	r *http.Request,
	err error,
) string {
	switch {
	case errors.Is(
		err,
		service.ErrCarInvalidID,
	):
		return translate(
			r,
			"car.invalid_id",
		)

	case errors.Is(
		err,
		service.ErrCarClientInvalidID,
	):
		return translate(
			r,
			"car.client_invalid_id",
		)

	case errors.Is(
		err,
		service.ErrCarOwnerInvalidID,
	):
		return translate(
			r,
			"car.owner_invalid_id",
		)

	case errors.Is(
		err,
		service.ErrCarBrandRequired,
	):
		return translate(
			r,
			"car.validation.brand_required",
		)

	case errors.Is(
		err,
		service.ErrCarModelRequired,
	):
		return translate(
			r,
			"car.validation.model_required",
		)

	case errors.Is(
		err,
		service.ErrCarPlateRequired,
	):
		return translate(
			r,
			"car.validation.plate_required",
		)

	case errors.Is(
		err,
		service.ErrCarClientRequired,
	):
		return translate(
			r,
			"car.validation.client_required",
		)

	case errors.Is(
		err,
		service.ErrCarAlreadyExists,
	):
		return translate(
			r,
			"car.already_exists",
		)

	default:
		return translate(
			r,
			"cars.error.internal",
		)
	}
}
