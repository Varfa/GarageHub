package handler

import "net/http"

func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "reports")
}
