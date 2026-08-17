package middleware

import (
	"context"
	"net/http"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type AuthMiddleware struct {
	sessionService *service.SessionService
	userService    *service.UserService
}
type contextKey string

const currentUserKey contextKey = "currentUser"

func NewAuthMiddleware(
	sessionService *service.SessionService,
	userService *service.UserService,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
		userService:    userService,
	}
}
func (m *AuthMiddleware) RequireAuth(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			cookie, err := r.Cookie("session")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			user, err := m.sessionService.GetUser(r.Context(), cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(
				r.Context(),
				currentUserKey,
				user,
			)
			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)

		},
	)
}
func CurrentUser(
	r *http.Request,
) (*models.User, bool) {
	user, ok := r.Context().
		Value(currentUserKey).(*models.User)

	return user, ok
}
func (m *AuthMiddleware) RequirePermission(
	permissionCode string,
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			user, ok := CurrentUser(r)

			if !ok {
				http.Error(
					w,
					http.StatusText(http.StatusUnauthorized),
					http.StatusUnauthorized,
				)
				return
			}

			if user.IsOwner {
				next.ServeHTTP(w, r)
				return
			}

			hasPermission, err := m.userService.HasPermission(
				r.Context(),
				user.ID,
				permissionCode,
			)
			if err != nil {
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}

			if !hasPermission {
				http.Error(
					w,
					http.StatusText(http.StatusForbidden),
					http.StatusForbidden,
				)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
