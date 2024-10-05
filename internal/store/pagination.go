package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
}

func (p *PaginatedFeedQuery) Parse(r *http.Request) error {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		p.Limit = limitInt
	}

	offset := qs.Get("offset")
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		p.Offset = offsetInt
	}

	sort := qs.Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		p.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		p.Search = search
	}

	return nil
}
