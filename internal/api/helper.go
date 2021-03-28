package api

import (
	"net/http"
	"strconv"
)

func findBegin(r *http.Request) int64 {
	pagesParam := r.URL.Query().Get("page")
	if pagesParam == "" {
		pagesParam = "0"
	}
	skip, err := strconv.ParseInt(pagesParam, 10, 64)
	if err != nil {
		return 0
	}
	return limit * skip
}
