package env

import (
	"context"
	"hub_project/models"
	"net/http"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
)

const (
	AuthAccountKey = "auth-account"
	SessionIDKey   = "SessionId"
)

func GetAuthAccount(r *http.Request) (a *models.Account) {
	a, ok := r.Context().Value(AuthAccountKey).(*models.Account)
	if !ok {
		panic("auth account not found in context of request")
	}
	return
}

func (e *env) Auth(next http.Handler) http.Handler {
	if e.sessionStore == nil {
		panic("can't provide auth middleware with nil session store")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r, err := e.sessionStore.Interpreter.PutObjType(r, &models.Account{})
		if err != nil {
			payload.ErrorResponse(w, err)
			return
		}

		session, err := e.sessionStore.Get(r, SessionIDKey)
		if err != nil {
			payload.ErrorResponse(w, err)
			return
		}

		account, isLoggedIn := e.GetFromSession(session)
		if !isLoggedIn {
			payload.ErrorResponse(w, payload.NewUnauthorizedErr("login before accessing the resource"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), AuthAccountKey, account))

		next.ServeHTTP(w, r)
		return
	})
}

func (e *env) AuthAdmin(next http.Handler) http.Handler {
	return e.Auth(e.isAdmin(next))
}
