package main

import (
	"net/http"
	"strconv"
	"errors"

	"github.com/biboyqg/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload

	if err := readJSON(w, r, &payload); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := 1

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: get user id from auth
		UserID:  int64(userID),
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNoRecord):
			errorJSON(w, http.StatusNotFound, err.Error())
		default:
			errorJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, &post); err != nil {
		errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}
