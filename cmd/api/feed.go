package main

import (
	"net/http"

	"github.com/biboyqg/social/internal/store"
)

//	@Summary		Get User Feed
//	@Description	Get the feed of a user by user ID
//	@Tags			Feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int			false	"Limit"
//	@Param			offset	query		int			false	"Offset"
//	@Param			sort	query		string		false	"Sort"
//	@Param			tags	query		[]string	false	"Tags"
//	@Param			search	query		string		false	"Search"
//	@Success		200		{array}		store.PostWithMetadata
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}

	if err := p.Parse(r); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(p); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.getUserFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	feed, err := app.store.Posts.GetUserFeed(ctx, user.ID, p)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
