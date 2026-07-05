package handler

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, page string) {
	tmpl, err := template.ParseFiles(
		"../frontend/templates/layout.html",
		"../frontend/templates/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
