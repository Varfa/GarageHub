package handler

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/service"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{
		service: service,
	}
}

func (h *ClientHandler) Clients(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "clients")
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {

}
