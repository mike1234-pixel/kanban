package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "session"
	sessionDuration   = 7 * 24 * time.Hour
)

// dummyHash is compared against when no user has the email, so a login for an
// unknown email takes as long as one with a wrong password. Otherwise the
// response time would reveal which emails have accounts.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

type AuthHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	user, err := h.Store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	hash := []byte(user.PasswordHash)
	if user.ID == "" {
		hash = dummyHash
	}
	// One message for both failures, so it never reveals whether the email exists.
	if err := bcrypt.CompareHashAndPassword(hash, []byte(payload.Password)); err != nil || user.ID == "" {
		writeError(w, http.StatusUnauthorized, codeInvalidCredentials, "Invalid email or password")
		return
	}

	token := newSessionToken()
	expiresAt := time.Now().Add(sessionDuration)
	if err := h.Store.CreateSession(r.Context(), hashToken(token), user.ID, expiresAt); err != nil {
		writeInternalError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token, // the real token goes only to the browser
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, user)
}

// #endregion

// #region Logout
// Logout deletes the session and clears the cookie. It always returns 204,
// even without a valid session, because the result is the same: logged out.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if err := h.Store.DeleteSession(r.Context(), hashToken(cookie.Value)); err != nil {
			writeInternalError(w, r, err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/", // must match the cookie's Path, or the browser won't delete it
		MaxAge:   -1,  // sends Max-Age=0: delete the cookie now
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

// #endregion

// #region Me
// Me returns the logged-in user. It must be wrapped in RequireAuth.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// #endregion

// newSessionToken returns 32 random bytes as URL-safe text, to use as a session ID.
func newSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b) // crypto/rand: unpredictable; never returns an error since Go 1.24
	return base64.RawURLEncoding.EncodeToString(b)
}

// hashToken returns the SHA-256 of a session token, as stored in sessions.token_hash.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
