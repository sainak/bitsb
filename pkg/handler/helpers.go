package handler

import (
	"net/http"
	"strconv"
)

const maxLimit = 50

// GetLimit returns the integer value of `limit` from the request
// if limit is not present in the request, it returns 10
// if limit is greater than maxLimit, it returns maxLimit
func GetLimit(r *http.Request) int64 {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit == 0 {
		limit = 10
	} else if limit > maxLimit {
		limit = maxLimit
	}
	return limit
}
