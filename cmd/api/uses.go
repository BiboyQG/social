package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/biboyqg/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userContextKey string
const userCtxKey userContextKey = "user"


func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtxKey).(*store.User)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNoRecord):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
