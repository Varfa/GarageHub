package handler

import (
	"net/http"
	"net/url"
	"strconv"
)

func (h *EmployeeHandler) SetPrimaryPhone(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "employees.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "employees.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	employeeID, err := strconv.ParseInt(
		r.FormValue("employee_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "employee.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	phoneID, err := strconv.ParseInt(
		r.FormValue("phone_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "employee.phone.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.SetPrimaryPhone(
		r.Context(),
		employeeID,
		phoneID,
	); err != nil {
		redirectEmployeePhoneError(
			w,
			r,
			employeeID,
			employeeErrorMessage(r, err),
		)
		return
	}

	redirectEmployeePhoneSuccess(
		w,
		r,
		employeeID,
		translate(r, "employee.success.primary_phone_changed"),
	)
}

func (h *EmployeeHandler) DeletePhone(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			translate(r, "employees.request.method_not_allowed"),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			translate(r, "employees.request.invalid"),
			http.StatusBadRequest,
		)
		return
	}

	employeeID, err := strconv.ParseInt(
		r.FormValue("employee_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "employee.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	phoneID, err := strconv.ParseInt(
		r.FormValue("phone_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			translate(r, "employee.phone.invalid_id"),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.DeletePhone(
		r.Context(),
		employeeID,
		phoneID,
	); err != nil {
		redirectEmployeePhoneError(
			w,
			r,
			employeeID,
			employeeErrorMessage(r, err),
		)
		return
	}

	redirectEmployeePhoneSuccess(
		w,
		r,
		employeeID,
		translate(r, "employee.success.phone_deleted"),
	)
}

func redirectEmployeePhoneError(
	w http.ResponseWriter,
	r *http.Request,
	employeeID int64,
	message string,
) {
	errorMessage := url.QueryEscape(
		message,
	)

	http.Redirect(
		w,
		r,
		"/employees/view?id="+
			strconv.FormatInt(employeeID, 10)+
			"&error="+errorMessage,
		http.StatusSeeOther,
	)
}

func redirectEmployeePhoneSuccess(
	w http.ResponseWriter,
	r *http.Request,
	employeeID int64,
	message string,
) {
	successMessage := url.QueryEscape(
		message,
	)

	http.Redirect(
		w,
		r,
		"/employees/view?id="+
			strconv.FormatInt(employeeID, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}
