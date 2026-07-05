package handler

import "net/http"

func EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "employees")
}
