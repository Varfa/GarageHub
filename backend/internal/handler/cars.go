package handler

import "net/http"

func CarsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "cars")
}
