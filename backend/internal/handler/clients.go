package handler

import "net/http"

func ClientsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "clients")
}
