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

//	@Summary		Get User
//	@Description	Get the user by user ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
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

//	@Summary		Follow User
//	@Description	Follow a user by user ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		201
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	follower, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: Revert back to auth userID from ctx
	userToFollowID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(ctx, userToFollowID, follower.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrAlreadyExists):
			app.conflict(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

//	@Summary		Unfollow User
//	@Description	Unfollow a user by user ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		200
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	follower, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	userToUnfollowID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(ctx, userToUnfollowID, follower.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

//	@Summary		Get Followers
//	@Description	Get the followers of a user by user ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{array}		store.User
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/followers [get]
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

//	@Summary		Activate User
//	@Description	Activate a user by user ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		201		{string}	string	"User activated"
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(ctx, token); err != nil {
		switch {
		case errors.Is(err, store.ErrNoRecord):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
