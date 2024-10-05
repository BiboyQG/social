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
	user, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
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

type authUserPayload struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userToFollow, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: Revert back to auth userID from ctx
	var payload authUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(ctx, userToFollow.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ErrAlreadyExists):
			app.conflict(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userToUnfollow, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: Revert back to auth userID from ctx
	var payload authUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(ctx, userToUnfollow.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	followers, err := app.store.Followers.GetFollowers(ctx, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, followers); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserFromCtx(r *http.Request) (*store.User, error) {
	ctx := r.Context()
	user, ok := ctx.Value(userCtxKey).(*store.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
