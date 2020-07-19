package env

import (
	"net/http"

	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"bitbucket.org/jazzserve/webapi/web/http/payload"
)

func (e *env) isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acc := GetAuthAccount(r)
		if val.Bool(acc.IsAdmin) == true {
			next.ServeHTTP(w, r)
			return
		}

		payload.ErrorResponse(w, payload.NewForbiddenErr("%s is not admin", val.Str(acc.Name)))
		return
	})
}
