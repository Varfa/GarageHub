package handler

import "net/http"

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "settings", nil)
}
