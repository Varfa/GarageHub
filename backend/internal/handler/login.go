package handler

import (
	"html/template"
	"net/http"

	"github.com/Varfa/GarageHub/internal/i18n"
)

type LoginHandler struct {
	translator *i18n.Manager
}
type LoginPageData struct {
	Language string
}

func NewLoginHandler(translator *i18n.Manager) *LoginHandler {
	return &LoginHandler{
		translator: translator,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	language := "en"

	cookie, err := r.Cookie("language")
	if err == nil && cookie.Value != "" {
		language = cookie.Value
	}

	funcMap := template.FuncMap{
		"t": func(key string) string {
			return h.translator.Translate(language, key)
		},
	}

	tmpl, err := template.New("login.html").
		Funcs(funcMap).
		ParseFiles("../frontend/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := LoginPageData{
		Language: language,
	}

	if err := tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
