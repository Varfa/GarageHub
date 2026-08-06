package handler

import "net/http"

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "orders", nil)
}
