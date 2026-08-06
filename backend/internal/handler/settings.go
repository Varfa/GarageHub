package handler

import "net/http"

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "settings", nil)
}
