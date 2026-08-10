package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Varfa/GarageHub/internal/i18n"
)

var translator *i18n.Manager

type TemplateData struct {
	PageData   any
	CurrentURL string
}

func SetTranslator(manager *i18n.Manager) {
	translator = manager
}

func requestLanguage(
	r *http.Request,
) string {
	language := "en"

	cookie, err := r.Cookie("language")
	if err == nil &&
		translator != nil &&
		translator.HasLanguage(cookie.Value) {
		language = cookie.Value
	}

	return language
}

func translate(
	r *http.Request,
	key string,
) string {
	if translator == nil {
		return key
	}

	return translator.Translate(
		requestLanguage(r),
		key,
	)
}

func RenderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	page string,
	data any,
) {
	language := requestLanguage(r)

	funcMap := template.FuncMap{
		"t": func(key string) string {
			return translator.Translate(
				language,
				key,
			)
		},

		"formatMoney": func(cents int64) string {
			euros := cents / 100
			remainder := cents % 100

			return fmt.Sprintf(
				"%d.%02d",
				euros,
				remainder,
			)
		},

		"lte": func(a, b int) bool {
			return a <= b
		},

		"positionName": func(
			code string,
			fallback string,
		) string {
			key := "employee.position." + code

			translated := translator.Translate(
				language,
				key,
			)

			if translated == key {
				return fallback
			}

			return translated
		},
		"phoneLabel": func(
			code string,
		) string {
			key := "employee.phones.type." + code

			translated := translator.Translate(
				language,
				key,
			)

			if translated == key {
				return code
			}

			return translated
		},
	}

	tmpl, err := template.New("layout.html").
		Funcs(funcMap).
		ParseFiles(
			"../frontend/templates/layout.html",
			"../frontend/templates/"+page+".html",
		)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	templateData := TemplateData{
		PageData:   data,
		CurrentURL: r.URL.RequestURI(),
	}

	err = tmpl.ExecuteTemplate(
		w,
		"layout.html",
		templateData,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}
