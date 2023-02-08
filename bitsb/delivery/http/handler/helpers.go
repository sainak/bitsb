package handler

import (
	"net/http"
	"strconv"
)

const maxLimit = 50

func getLimit(r *http.Request) int64 {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit == 0 {
		limit = 10
	} else if limit > maxLimit {
		limit = maxLimit
	}
	return limit
}
