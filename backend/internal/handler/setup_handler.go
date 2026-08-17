package handler

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/service"
)

type SetupPageData struct {
	Error string
}
type SetupHandler struct {
	userService *service.UserService
}

func NewSetupHandler(
	userService *service.UserService,
) *SetupHandler {
	return &SetupHandler{
		userService: userService,
	}
}

func (h *SetupHandler) SetupPage(
	w http.ResponseWriter,
	r *http.Request,
) {
	hasOwner, err := h.userService.HasOwner(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"Internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	if hasOwner {
		http.Redirect(
			w,
			r,
			"/login",
			http.StatusSeeOther,
		)
		return
	}

	data := SetupPageData{}

	RenderTemplate(w, r, "setup", data)

}
func (h *SetupHandler) CreateOwner(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	if password != passwordConfirm {
		data := SetupPageData{
			Error: translate(
				r,
				"setup.password_mismatch",
			),
		}

		RenderTemplate(
			w,
			r,
			"setup",
			data,
		)
		return
	}

	err = h.userService.CreateOwner(r.Context(), email, password)

	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
