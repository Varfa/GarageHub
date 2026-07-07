package handler

import (
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type ClientHandler struct {
	service *service.ClientService
}
type ClientsPageData struct {
	Clients []models.Client
	Error   string
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{
		service: service,
	}
}

func (h *ClientHandler) Clients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := ClientsPageData{
		Clients: clients,
	}
	RenderTemplate(w, "clients", data)

}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	client := models.Client{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		Address: r.FormValue("address"),
		Note:    r.FormValue("note"),
	}
	err = h.service.Create(r.Context(), client)

	if err != nil {
		clients, listErr := h.service.List(r.Context())
		if listErr != nil {
			http.Error(w, listErr.Error(), http.StatusInternalServerError)
			return
		}

		data := ClientsPageData{
			Clients: clients,
			Error:   err.Error(),
		}

		RenderTemplate(w, "clients", data)
		return
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

//УДАЛЕНИЕ КЛИЕНТА ИЗ БД

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}
