package middleware

import (
	"context"
	"net/http"
	"time"

	"todo-list-app/db"
	"todo-list-app/model"
)

func isExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			ShowMessage(w, "You are not Logged in!", 401)
			return
		}
		sessionToken := c.Value
		userSession, exists := db.Sessions[sessionToken]

		if !exists {
			ShowMessage(w, "You are not Logged in!", 401)
			return
		}

		if isExpired(userSession) {
			delete(db.Sessions, sessionToken)
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Path:    "/",
				Expires: time.Now(),
			})
			ShowMessage(w, "You are not Logged in!", 401)
			return
		}

		ctx := context.WithValue(r.Context(), "username", userSession.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}) // TODO: replace this
}
