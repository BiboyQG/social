package main

import (
	"errors"
	"net/http"
	"strconv"
	"context"

	"github.com/biboyqg/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postCtxKey string
const postContextKey postCtxKey = "post"

type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type updatePostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    []string  `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	userID := 1

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: get user id from auth
		UserID: int64(userID),
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, &post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	post, err := app.getPostFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	post, err := app.getPostFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload updatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.Update(ctx, post); err != nil {
		switch {
		case errors.Is(err, store.ErrNoRecord):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	post, err := app.getPostFromCtx(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNoRecord):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{"message": "post deleted"}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNoRecord):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getPostFromCtx(r *http.Request) (*store.Post, error) {
	ctx := r.Context()
	post, ok := ctx.Value(postContextKey).(*store.Post)
	if !ok {
		return nil, errors.New("post not found in context")
	}
	return post, nil
}
