package auth

import (
	"log"
	"net/http"

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
		userEmail, emailExists := session.Values["user_email"].(string)

		log.Printf("AuthRequired Middleware: authenticated=%v, userEmail=%v", authenticated, userEmail)

		if !ok || !authenticated {
			log.Println("User is not authenticated. Redirecting to login.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if emailExists {
			log.Printf("Authenticated user: %s", userEmail)
		}

		next.ServeHTTP(w, r)
	})
}
