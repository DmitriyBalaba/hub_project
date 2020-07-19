package handlers

import (
	"context"
	"hub_project/env"
	"hub_project/handlers/dto"
	"hub_project/models"
	"net/http"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
	"bitbucket.org/jazzserve/webapi/web/http/routing"
	"github.com/gorilla/mux"
)

type handlerAStores struct {
	env.Environment
}

func init() {
	routing.AddRouteInitializer(func(r *mux.Router, ctx context.Context) {
		h := &handlerAStores{env.Get(ctx)}

		r = r.PathPrefix("/account-stores").Subrouter()

		routing.New(r, "", h.postStore, h.Auth,
			payload.JsonBody(dto.AccountStorePost{}, models.AccountStore{})).
			Methods(http.MethodPost)
	})
}

func (h *handlerAStores) postStore(w http.ResponseWriter, r *http.Request) (err error) {
	model := payload.GetModel(r.Context()).(*models.AccountStore)

	err = h.Storage().CreateAccountStore(model)

	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	return nil
}
