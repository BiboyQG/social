package main

import (
	"net/http"

	"github.com/biboyqg/social/internal/store"
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
