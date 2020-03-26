package main

import (
	"context"
	"net/http"

	goctx "github.com/gorilla/context"
)

// BasicUserCheck doing a simple `uid` validation if exists query params.
// example: if endpoint accessed with `?uid=abc123` then that uid will be checked
// to see if its exists on DB, if exists then allow the request, otherwise returns bad request error.
func (s *server) BasicUserCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		uid := query.Get("uid")

		if uid == "" {
			next.ServeHTTP(w, r)
		} else {
			ctx := context.Background()
			auth, err := s.app.Auth(ctx)
			if err != nil {
				http.Error(w, "auth server error", http.StatusInternalServerError)
				return
			}
			// TODO: to add more security we could verify the ID tokens here instead of just checking
			// if uid is associated with a registered user.
			user, err := auth.GetUser(ctx, uid)
			if err != nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			// pass user object to the next request
			goctx.Set(r, userCtxKey, user)
			next.ServeHTTP(w, r)
		}
	})
}
