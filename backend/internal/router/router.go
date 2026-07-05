package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	static := http.FileServer(http.Dir("../frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))
	mux.HandleFunc("/", handler.HomeHandler)

	mux.HandleFunc("/dashboard", handler.DashboardHandler)

	mux.HandleFunc("/health", handler.HealthHandler)

	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API 👋"))
	})

	return mux

}
