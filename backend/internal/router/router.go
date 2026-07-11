package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
)

func SetupRoutes(clientHandler *handler.ClientHandler,
	carHandler *handler.CarHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	static := http.FileServer(http.Dir("../frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	// Pages
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/dashboard", handler.DashboardHandler)

	// CLIENTS

	mux.HandleFunc("/clients", clientHandler.Clients)
	mux.HandleFunc("/clients/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			clientHandler.CreatePage(w, r)
			return
		}

		clientHandler.Create(w, r)
	})
	mux.HandleFunc("/clients/delete", clientHandler.Delete)
	mux.HandleFunc("/clients/view", clientHandler.View)
	mux.HandleFunc("/clients/update", clientHandler.Update)
	mux.HandleFunc("/orders", handler.OrdersHandler)
	mux.HandleFunc("/parts", handler.PartsHandler)
	mux.HandleFunc("/warehouse", handler.WarehouseHandler)
	mux.HandleFunc("/employees", handler.EmployeesHandler)
	mux.HandleFunc("/reports", handler.ReportsHandler)
	mux.HandleFunc("/settings", handler.SettingsHandler)

	// CARS
	mux.HandleFunc("/cars/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			carHandler.CreatePage(w, r)
			return
		}

		carHandler.Create(w, r)
	})
	mux.HandleFunc("/cars/delete", carHandler.Delete)
	mux.HandleFunc("/cars", carHandler.Cars)
	mux.HandleFunc("/cars/view", carHandler.View)
	mux.HandleFunc("/cars/update", carHandler.Update)

	// System
	mux.HandleFunc("/health", handler.HealthHandler)

	// API
	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API 👋"))
	})

	return mux
}
