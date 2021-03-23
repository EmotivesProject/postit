package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"postit/model"
)

type key string

const (
	userID key = "username"
)

var (
	errUnauthorised = errors.New("Not authorised")
)

func verifyJTW() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := &http.Client{}

			authURL := "http://uacl/authorize"

			header := r.Header.Get("Authorization")

			req, err := http.NewRequest("GET", authURL, nil)
			if err != nil {
				messageResponseJSON(w, http.StatusBadRequest, model.Message{
					Message: err.Error(),
				})
				return
			}
			req.Header.Add("Authorization", header)
			resp, err := client.Do(req)
			if err != nil {
				messageResponseJSON(w, http.StatusBadRequest, model.Message{
					Message: err.Error(),
				})
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					messageResponseJSON(w, http.StatusBadRequest, model.Message{
						Message: err.Error(),
					})
					return
				}

				var dat model.Response
				var user model.User
				_ = json.Unmarshal(bodyBytes, &dat)
				body, _ := json.Marshal(dat.Result)
				_ = json.Unmarshal(body, &user)
				ctx := context.WithValue(r.Context(), userID, user.Username)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			messageResponseJSON(w, http.StatusBadRequest, model.Message{
				Message: errUnauthorised.Error(),
			})
		})
	}
}

func SimpleMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*") // fixme please
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
