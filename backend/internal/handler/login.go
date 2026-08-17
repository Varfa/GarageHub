package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/Varfa/GarageHub/internal/i18n"
	"github.com/Varfa/GarageHub/internal/service"
)

type LoginHandler struct {
	translator     *i18n.Manager
	userService    *service.UserService
	sessionService *service.SessionService
}
type LoginPageData struct {
	Language string
	Error    string
	Email    string
}

func NewLoginHandler(
	translator *i18n.Manager,
	userService *service.UserService,
	sessionService *service.SessionService,
) *LoginHandler {
	return &LoginHandler{
		translator:     translator,
		userService:    userService,
		sessionService: sessionService,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
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
		remember := r.FormValue("remember") == "on"

		user, err := h.userService.Authenticate(
			r.Context(),
			email,
			password,
		)
		if err != nil {
			if errors.Is(
				err,
				service.ErrInvalidCredentials,
			) {
				h.renderLogin(
					w,
					r,
					translate(r, "login.invalid_credentials"),
					email,
				)
				return
			}
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		token, expiresAt, err := h.sessionService.Create(
			r.Context(),
			user.ID,
			remember,
		)
		if err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(
			w,
			r,
			"/dashboard",
			http.StatusSeeOther,
		)
		return

	}
	h.renderLogin(
		w,
		r,
		"",
		"",
	)
}
func (h *LoginHandler) renderLogin(
	w http.ResponseWriter,
	r *http.Request,
	errorMessage string,
	email string,
) {
	language := "en"

	cookie, err := r.Cookie("language")
	if err == nil && cookie.Value != "" {
		language = cookie.Value
	}

	funcMap := template.FuncMap{
		"t": func(key string) string {
			return h.translator.Translate(
				language,
				key,
			)
		},
	}

	tmpl, err := template.New("login.html").
		Funcs(funcMap).
		ParseFiles("../frontend/templates/login.html")
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := LoginPageData{
		Language: language,
		Error:    errorMessage,
		Email:    email,
	}

	if err := tmpl.ExecuteTemplate(
		w,
		"login.html",
		data,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}

func RootRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
func (h *LoginHandler) ChangeLanguage(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("lang")

	if !h.translator.HasLanguage(language) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	redirectTo := r.URL.Query().Get("redirect")
	if redirectTo == "" || redirectTo[0] != '/' {
		redirectTo = "/dashboard"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "language",
		Value:    language,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}
func (h *LoginHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		err = h.sessionService.Delete(r.Context(), cookie.Value)
		if err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)

	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)

}
