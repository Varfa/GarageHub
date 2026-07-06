package handler

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, page string, data any) {
	tmpl, err := template.ParseFiles(
		"../frontend/templates/layout.html",
		"../frontend/templates/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
