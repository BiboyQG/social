package main

import (
	"net/http"

	"github.com/biboyqg/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	if err := p.Parse(r); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(p); err != nil {
		app.badRequest(w, r, err)
		return
	}

	userID := int64(1)
	// userID := ctx.Value(userIDKey).(int64)

	feed, err := app.store.Posts.GetUserFeed(ctx, userID, p)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
