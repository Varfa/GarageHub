package middleware

import (
	"net/http"
)

func RequireOwner(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			user, ok := CurrentUser(r)

			if !ok || user == nil {
				http.Error(
					w,
					http.StatusText(http.StatusUnauthorized),
					http.StatusUnauthorized,
				)
				return
			}

			if !user.IsOwner {
				http.Error(
					w,
					http.StatusText(http.StatusForbidden),
					http.StatusForbidden,
				)
				return
			}

			next.ServeHTTP(
				w,
				r,
			)
		},
	)
}
