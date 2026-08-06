package handler

import (
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "dashboard", nil)

}
