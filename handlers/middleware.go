package handlers

import (
	"app/services"
	"net/http"
)

type Middleware struct {
}

func (m *Middleware) HasAccess(au *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spaceId := r.PathValue("id")
		sessionCookie, err := r.Cookie("session")

		var sessionString string

		if err != nil {
			sessionString = ""
		} else {
			sessionString = sessionCookie.Value
		}

		if au.HasAccess(sessionString, spaceId) {
			next(w, r)
		} else {
			http.Error(w, "unauthorized", http.StatusForbidden)
		}
	}
}
