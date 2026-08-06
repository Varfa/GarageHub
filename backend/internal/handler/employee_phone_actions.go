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
			"метод не поддерживается",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
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
			"некорректный id сотрудника",
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
			"некорректный id телефона",
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
			err,
		)
		return
	}

	redirectEmployeePhoneSuccess(
		w,
		r,
		employeeID,
		"основной телефон изменён",
	)
}

func (h *EmployeeHandler) DeletePhone(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"метод не поддерживается",
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"некорректный запрос",
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
			"некорректный id сотрудника",
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
			"некорректный id телефона",
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
			err,
		)
		return
	}

	redirectEmployeePhoneSuccess(
		w,
		r,
		employeeID,
		"телефон успешно удалён",
	)
}

func redirectEmployeePhoneError(
	w http.ResponseWriter,
	r *http.Request,
	employeeID int64,
	err error,
) {
	errorMessage := url.QueryEscape(
		err.Error(),
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
	successMessage := url.QueryEscape(message)

	http.Redirect(
		w,
		r,
		"/employees/view?id="+
			strconv.FormatInt(employeeID, 10)+
			"&success="+successMessage,
		http.StatusSeeOther,
	)
}
