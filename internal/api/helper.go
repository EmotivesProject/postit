package api

import (
	"net/http"
	"postit/internal/db"
	"strconv"

	"github.com/go-chi/chi"
)

func findBegin(r *http.Request) int {
	pagesParam := r.URL.Query().Get("page")
	if pagesParam == "" {
		pagesParam = "0"
	}

	skip, err := strconv.Atoi(pagesParam)
	if err != nil {
		return 0
	}

	return db.PostLimit * skip
}

func extractID(r *http.Request, param string) (int, error) {
	paramString := chi.URLParam(r, param)

	return strconv.Atoi(paramString)
}
