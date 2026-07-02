package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home page 🚀"))
	})

	mux.HandleFunc("/health", handler.HealthHandler)

	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API 👋"))
	})

	return mux
}
