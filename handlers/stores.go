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

type handlerStores struct {
	env.Environment
}

func init() {
	routing.AddRouteInitializer(func(r *mux.Router, ctx context.Context) {
		h := &handlerStores{env.Get(ctx)}

		r = r.PathPrefix("/stores").Subrouter()

		routing.New(r, "", h.getStoreList, h.Auth,
			payload.QueryParams(payload.OffsetLimit{})).
			Methods(http.MethodGet)

		routing.New(r, "", h.postStore, h.Auth,
			payload.JsonBody(dto.StorePost{}, dto.StorePost{})).
			Methods(http.MethodPost)

		routing.New(r, "/{id}", h.putStore, h.Auth,
			payload.JsonBody(dto.StorePut{}, dto.StorePut{})).
			Methods(http.MethodPut)

		routing.New(r, "/{id}", h.deleteStore, h.Auth).
			Methods(http.MethodDelete)

		routing.New(r, "/accounts", h.getStoreAccountList, h.Auth,
			payload.QueryParams(dto.AccountStoreQuery{})).
			Methods(http.MethodGet)
	})
}

func (h *handlerStores) getStoreList(w http.ResponseWriter, r *http.Request) (err error) {
	reqDto := payload.GetQueryParams(r.Context()).(*payload.OffsetLimit)
	offset, limit := reqDto.GetOffsetLimit()

	out := make([]*models.Store, 0, limit)
	err = h.Storage().
		Preload("AccountStoreList.Account").
		LeftJoin("account_store a ON a.store_id = store.id").
		OffsetAndLimit(offset, limit).
		SortAsc("a.account_id").
		GetListOf(&out)
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	payload.JsonEncode(w, out)
	return nil
}

func (h *handlerStores) getStoreAccountList(w http.ResponseWriter, r *http.Request) (err error) {
	reqDto := payload.GetQueryParams(r.Context()).(*dto.AccountStoreQuery)
	offset, limit := reqDto.GetOffsetLimit()

	out := make([]*models.Store, 0, limit)
	err = h.Storage().
		PreloadAccount(reqDto.AccountID).
		Preload("AccountStoreList.Store").
		LeftJoin("account_store a ON a.store_id = store.id").
		OffsetAndLimit(offset, limit).
		SortAsc("a.account_id").
		GetListOf(&out)
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	payload.JsonEncode(w, out)
	return nil
}

func (h *handlerStores) postStore(w http.ResponseWriter, r *http.Request) (err error) {
	reqDto := payload.GetDTO(r.Context()).(*dto.StorePost)

	in := &models.Store{
		Name: reqDto.Name,
	}
	err = h.Storage().Create(in)

	if err != nil {
		return err
	}
	payload.JsonEncode(w, in)
	return nil
}

func (h *handlerStores) putStore(w http.ResponseWriter, r *http.Request) (err error) {
	ID, err := payload.GetPathInt64(r, "id")
	if err != nil {
		return payload.WrapBadRequestErr(err)
	}

	reqDto := payload.GetDTO(r.Context()).(*dto.StorePut)

	in := &models.Store{
		ID:               &ID,
		Name:             reqDto.Name,
		AccountStoreList: reqDto.AccountStoreList,
	}
	err = h.Storage().UpdateStore(in)
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	in.ID = &ID
	payload.JsonEncode(w, in)
	return nil
}

func (h *handlerStores) deleteStore(w http.ResponseWriter, r *http.Request) (err error) {
	ID, err := payload.GetPathInt64(r, "id")
	if err != nil {
		return payload.WrapBadRequestErr(err)
	}

	return h.Storage().DeleteStore(&ID)
}
