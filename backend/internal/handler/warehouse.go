package handler

import "net/http"

func WarehouseHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "warehouse")
}
