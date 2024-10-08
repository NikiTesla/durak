package api

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	usernameKey ContextKey = "username"
	roleKey     ContextKey = "admin"

	adminRole  = "admin"
	playerRole = "player"
)

func (a *App) playerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Debug("player middleware is serving request")

		claims, err := parseClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
		ctx = context.WithValue(ctx, roleKey, playerRole)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *App) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Debug("admin middleware is serving request")

		claims, err := parseClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if getRole(claims.Username) != adminRole {
			http.Error(w, "you are not allowed to do that", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
		ctx = context.WithValue(ctx, roleKey, adminRole)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
