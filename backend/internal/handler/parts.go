package handler

import "net/http"

func PartsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "parts", nil)
}
