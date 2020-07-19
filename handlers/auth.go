package handlers

import (
	"context"
	"hub_project/env"
	"hub_project/handlers/dto"
	"hub_project/models"
	"net/http"
	"strings"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
	"bitbucket.org/jazzserve/webapi/web/http/routing"
	"github.com/gorilla/mux"
)

type handlerAbout struct {
	env.Environment
}

func init() {
	routing.AddRouteInitializer(func(r *mux.Router, ctx context.Context) {
		h := &handlerAbout{env.Get(ctx)}

		r = r.PathPrefix("/").Subrouter()

		routing.New(r, "/login", h.postLogin,
			payload.JsonBody(dto.AccountLogin{}, dto.AccountLogin{})).
			Methods(http.MethodPost)

		routing.New(r, "/logout", h.postLogout).
			Methods(http.MethodPost)
	})
}

func (h *handlerAbout) postLogin(w http.ResponseWriter, r *http.Request) (err error) {
	credentials := payload.GetDTO(r.Context()).(*dto.AccountLogin)
	credentials.Email = strings.TrimSpace(credentials.Email)
	credentials.Password = strings.TrimSpace(credentials.Password)

	session, r, err := h.GetSession(r)
	if err != nil {
		return payload.WrapUnauthorizedErr(err)
	}

	account, isLoggedIn := h.GetFromSession(session)
	if isLoggedIn && account.HasEmail(credentials.Email) {
		payload.JsonEncode(w, account)
		return
	}

	account = &models.Account{}
	isFound, err := h.Storage().
		MatchEmail(credentials.Email).
		GetAnyOfOptional(account)
	if err != nil {
		return
	}

	if !isFound {
		return payload.NewUnauthorizedErr("email not found")
	}

	if !account.HasPassword(credentials.Password) {
		return payload.NewUnauthorizedErr("password is invalid")
	}

	if err = h.AddToSession(account, session); err != nil {
		return
	}

	if err = h.SaveSessionAsNew(r, w, session, account.ID); err != nil {
		return
	}

	payload.JsonEncode(w, account)
	return nil
}

func (h *handlerAbout) postLogout(w http.ResponseWriter, r *http.Request) (err error) {
	session, r, err := h.GetSession(r)
	if err != nil {
		return
	}

	_, isLoggedIn := h.GetFromSession(session)
	if !isLoggedIn {
		return
	}

	if err = h.DeleteSession(r, w, session); err != nil {
		return
	}
	return nil
}
