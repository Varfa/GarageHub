package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
)

func SetupRoutes(clientHandler *handler.ClientHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	static := http.FileServer(http.Dir("../frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	// Pages
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/dashboard", handler.DashboardHandler)
	mux.HandleFunc("/clients", clientHandler.Clients)
	mux.HandleFunc("/clients/create", clientHandler.Create)
	mux.HandleFunc("/clients/delete", clientHandler.Delete)
	mux.HandleFunc("/cars", handler.CarsHandler)
	mux.HandleFunc("/orders", handler.OrdersHandler)
	mux.HandleFunc("/parts", handler.PartsHandler)
	mux.HandleFunc("/warehouse", handler.WarehouseHandler)
	mux.HandleFunc("/employees", handler.EmployeesHandler)
	mux.HandleFunc("/reports", handler.ReportsHandler)
	mux.HandleFunc("/settings", handler.SettingsHandler)

	// System
	mux.HandleFunc("/health", handler.HealthHandler)

	// API
	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API 👋"))
	})

	return mux
}
