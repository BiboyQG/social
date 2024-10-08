package main

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorized(w, r, errors.New("authorization header is missing"))
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			app.unauthorized(w, r, errors.New("invalid authorization scheme"))
			return
		}

		token := authParts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorized(w, r, err)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !jwtToken.Valid {
			app.unauthorized(w, r, errors.New("invalid token claims"))
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			app.unauthorized(w, r, errors.New("invalid user ID in token claims"))
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, int64(userID))
		if err != nil {
			app.unauthorized(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the request header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicAuth(w, r, errors.New("authorization header is missing"))
				return
			}

			// Split the header into the scheme and the credentials
			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 || authParts[0] != "Basic" {
				app.unauthorizedBasicAuth(w, r, errors.New("invalid authorization scheme"))
				return
			}

			// Decode the credentials
			decodedCredentials, err := base64.StdEncoding.DecodeString(authParts[1])
			if err != nil {
				app.unauthorizedBasicAuth(w, r, errors.New("error while decoding credentials"))
				return
			}

			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			// Split the credentials into the username and password
			credentials := strings.SplitN(string(decodedCredentials), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
				app.unauthorizedBasicAuth(w, r, errors.New("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
