package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
)

func SetupRoutes(
	clientHandler *handler.ClientHandler,
	carHandler *handler.CarHandler,
	employeeHandler *handler.EmployeeHandler,
	loginHandler *handler.LoginHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	static := http.FileServer(
		http.Dir("../frontend/static"),
	)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", static),
	)

	// Pages
	mux.HandleFunc(
		"/language",
		loginHandler.ChangeLanguage,
	)

	mux.HandleFunc(
		"/",
		handler.RootRedirectHandler,
	)

	mux.HandleFunc(
		"/login",
		loginHandler.Login,
	)

	mux.HandleFunc(
		"/dashboard",
		handler.DashboardHandler,
	)

	// Clients
	mux.HandleFunc(
		"/clients",
		clientHandler.Clients,
	)

	mux.HandleFunc(
		"/clients/create",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				clientHandler.CreatePage(w, r)
				return
			}

			clientHandler.Create(w, r)
		},
	)

	mux.HandleFunc(
		"/clients/delete",
		clientHandler.Delete,
	)

	mux.HandleFunc(
		"/clients/view",
		clientHandler.View,
	)

	mux.HandleFunc(
		"/clients/update",
		clientHandler.Update,
	)

	// Cars
	mux.HandleFunc(
		"/cars",
		carHandler.Cars,
	)

	mux.HandleFunc(
		"/cars/create",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				carHandler.CreatePage(w, r)
				return
			}

			carHandler.Create(w, r)
		},
	)

	mux.HandleFunc(
		"/cars/delete",
		carHandler.Delete,
	)

	mux.HandleFunc(
		"/cars/view",
		carHandler.View,
	)

	mux.HandleFunc(
		"/cars/update",
		carHandler.Update,
	)

	mux.HandleFunc(
		"/cars/change-owner",
		carHandler.ChangeOwner,
	)

	// Employees
	mux.HandleFunc(
		"/employees",
		employeeHandler.Employees,
	)

	mux.HandleFunc(
		"/employees/archive",
		employeeHandler.ArchivePage,
	)

	mux.HandleFunc(
		"/employees/view",
		employeeHandler.View,
	)

	mux.HandleFunc(
		"/employees/update",
		employeeHandler.Update,
	)

	mux.HandleFunc(
		"/employees/archive-action",
		employeeHandler.Archive,
	)

	mux.HandleFunc(
		"/employees/restore",
		employeeHandler.Restore,
	)

	mux.HandleFunc(
		"/employees/create",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				employeeHandler.CreatePage(w, r)
				return
			}

			employeeHandler.Create(w, r)
		},
	)

	mux.HandleFunc(
		"/employees/phones/create",
		employeeHandler.AddPhone,
	)

	mux.HandleFunc(
		"/employees/phones/primary",
		employeeHandler.SetPrimaryPhone,
	)

	mux.HandleFunc(
		"/employees/phones/delete",
		employeeHandler.DeletePhone,
	)

	// Modules
	mux.HandleFunc(
		"/orders",
		handler.OrdersHandler,
	)

	mux.HandleFunc(
		"/parts",
		handler.PartsHandler,
	)

	mux.HandleFunc(
		"/warehouse",
		handler.WarehouseHandler,
	)

	mux.HandleFunc(
		"/reports",
		handler.ReportsHandler,
	)

	mux.HandleFunc(
		"/settings",
		handler.SettingsHandler,
	)

	// System
	mux.HandleFunc(
		"/health",
		handler.HealthHandler,
	)

	// API
	mux.HandleFunc(
		"/api/v1/hello",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(
				[]byte("Hello from API 👋"),
			)
		},
	)

	return mux
}
