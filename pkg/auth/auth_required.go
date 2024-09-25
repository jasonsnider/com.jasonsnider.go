package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/boj/redistore"
)

type AuthMiddleware struct {
	SessionStore *redistore.RediStore
}

func (m *AuthMiddleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.SessionStore.Get(r, "com-jasonsnider-go")
		if err != nil {
			log.Printf("Failed to get session: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Log session values to debug
		authenticated, ok := session.Values["authenticated"].(bool)

		if !ok || !authenticated {
			log.Println("User is not authenticated. Redirecting to login.")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		sessionExpiryStr := os.Getenv("SESSION_EXPIRY")
		sessionExpiry, err := strconv.Atoi(sessionExpiryStr)

		if err != nil {
			log.Printf("Invalid session expiry value: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Reset session expiration
		session.Options.MaxAge = sessionExpiry

		// Save the session to update the expiration time
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Failed to save session during AuthRequired: %v", err)
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
