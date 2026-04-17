package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"go-huginn-clone/models"
)

var store = sessions.NewCookieStore([]byte("huginn-secret-key-change-in-production"))

type contextKey string

const (
	CurrentUserKey contextKey = "current_user"
	FlashKey       contextKey = "flash"
)

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
}

// GetSession returns the gorilla session
func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "huginn-session")
	return session
}

// SaveSession saves the session
func SaveSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	session.Save(r, w)
}

// SetCurrentUser stores the mock current user in context
func SetCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := GetSession(r)
		// Check if the user is "logged in" via session flag
		loggedIn, _ := session.Values["logged_in"].(bool)
		var user *models.User
		if loggedIn {
			user = models.MockUser()
		}
		ctx := context.WithValue(r.Context(), CurrentUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCurrentUser extracts the current user from context
func GetCurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(CurrentUserKey).(*models.User)
	return user
}

// RequireLogin redirects to login if not authenticated
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetCurrentUser(r)
		if user == nil {
			http.Redirect(w, r, "/users/sign_in", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin redirects to root if not admin
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetCurrentUser(r)
		if user == nil || !user.Admin {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SetFlash stores a flash message in session
func SetFlash(w http.ResponseWriter, r *http.Request, msgType, message string) {
	session := GetSession(r)
	session.AddFlash(message, msgType)
	session.Save(r, w)
}

// GetFlashes retrieves and clears flash messages
func GetFlashes(w http.ResponseWriter, r *http.Request) []models.FlashMessage {
	session := GetSession(r)
	var flashes []models.FlashMessage
	for _, t := range []string{"notice", "alert"} {
		for _, f := range session.Flashes(t) {
			if msg, ok := f.(string); ok {
				flashes = append(flashes, models.FlashMessage{Type: t, Message: msg})
			}
		}
	}
	session.Save(r, w)
	return flashes
}

// LoginUser marks user as logged in via session
func LoginUser(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	session.Values["logged_in"] = true
	session.Save(r, w)
}

// LogoutUser clears the session
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	session.Values["logged_in"] = false
	delete(session.Values, "logged_in")
	session.Save(r, w)
}

// PageProps builds the layouts.PageProps for a request
func PagePropsFromRequest(w http.ResponseWriter, r *http.Request, title string) interface{} {
	// Returned as interface{} to avoid import cycle; handlers cast it
	return map[string]interface{}{
		"title":       title,
		"currentUser": GetCurrentUser(r),
		"flashes":     GetFlashes(w, r),
		"currentPath": r.URL.Path,
	}
}
