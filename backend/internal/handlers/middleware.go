package handlers

import (
	"context"
	"net/http"

	"kanban/internal/models"
)

// contextKey is unexported, so no other package can create the same key and
// overwrite (or read) the values this package stores in a request's context.
type contextKey struct{}

var userKey = contextKey{}

// RequireAuth only lets a request through to next if it has a valid session
// cookie. The logged-in user is added to the request's context.
func (h *AuthHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, codeUnauthenticated, "You must be logged in")
			return
		}

		user, err := h.Store.GetSessionUser(r.Context(), hashToken(cookie.Value))
		if err != nil {
			writeInternalError(w, r, err)
			return
		}
		// Unknown and expired sessions both come back as an empty user.
		if user.ID == "" {
			writeError(w, http.StatusUnauthorized, codeUnauthenticated, "You must be logged in")
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// userFromContext returns the user RequireAuth added to the context.
// ok is false if the request didn't go through RequireAuth.
func userFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userKey).(models.User)
	return user, ok
}
