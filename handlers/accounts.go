package handlers

import (
	"context"
	"fmt"
	"hub_project/env"
	"hub_project/handlers/dto"
	"hub_project/models"
	"net/http"
	"strings"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
	"bitbucket.org/jazzserve/webapi/web/http/routing"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type handlerUserManagement struct {
	env.Environment
}

func init() {
	routing.AddRouteInitializer(func(r *mux.Router, ctx context.Context) {
		h := &handlerUserManagement{env.Get(ctx)}

		r = r.PathPrefix("/accounts").Subrouter()

		routing.New(r, "", h.getAccountList, h.AuthAdmin,
			payload.QueryParams(dto.AccountGet{})).
			Methods(http.MethodGet)

		routing.New(r, "", h.postAccount, h.AuthAdmin,
			payload.JsonBody(dto.AccountPost{}, dto.AccountPost{})).
			Methods(http.MethodPost)

		routing.New(r, "/{id}", h.putAccount, h.AuthAdmin,
			payload.JsonBody(dto.AccountPut{}, models.Account{})).
			Methods(http.MethodPut)

		routing.New(r, "/{id}", h.deleteAccount, h.AuthAdmin).
			Methods(http.MethodDelete)
	})
}

func (h *handlerUserManagement) getAccountList(w http.ResponseWriter, r *http.Request) (err error) {
	reqDto := payload.GetQueryParams(r.Context()).(*dto.AccountGet)
	offset, limit := reqDto.GetOffsetLimit()

	out := make([]*models.Account, 0, limit)
	err = h.Storage().
		OffsetAndLimit(offset, limit).
		SortAsc("name").
		GetListOf(&out)
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	payload.JsonEncode(w, out)
	return nil
}

func (h *handlerUserManagement) postAccount(w http.ResponseWriter, r *http.Request) (err error) {
	credentials := payload.GetDTO(r.Context()).(*dto.AccountPost)
	credentials.Email = strings.TrimSpace(credentials.Email)

	ok, err := models.IsAccountExists(&credentials.Email, nil, h.Storage())
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	//TODO: add WrapStatusConflictErr
	if !ok {
		err = fmt.Errorf("user with email %s already exists", credentials.Email)
		return payload.WrapBadRequestErr(err)
	}

	acc := &models.Account{
		Name:           &credentials.Name,
		Email:          &credentials.Email,
		IsAdmin:        &credentials.IsAdmin,
		IsStoreManager: &credentials.IsStoreManager,
		AccountStore: &models.AccountStore{
			Store: credentials.Store,
		},
	}
	if err := acc.Create(h.Storage()); err != nil {
		return payload.WrapInternalServerErr(err)
	}

	if err := h.SendNotificationToNewAccount(acc); err != nil {
		log.Error().Msgf("cannot notification account [%s]", err.Error())
	}
	payload.JsonEncode(w, credentials)
	return nil
}

func (h *handlerUserManagement) putAccount(w http.ResponseWriter, r *http.Request) (err error) {
	ID, err := payload.GetPathInt64(r, "id")
	if err != nil {
		return payload.WrapBadRequestErr(err)
	}

	model := payload.GetModel(r.Context()).(*models.Account)

	ok, err := models.IsAccountExists(model.Email, &ID, h.Storage())
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	if !ok {
		return payload.NewForbiddenErr("account already exists")
	}

	err = h.Storage().Match("id", &ID).Update(model)
	if err != nil {
		return payload.WrapInternalServerErr(err)
	}

	model.ID = &ID
	payload.JsonEncode(w, model)
	return nil
}

func (h *handlerUserManagement) deleteAccount(w http.ResponseWriter, r *http.Request) (err error) {
	ID, err := payload.GetPathInt64(r, "id")
	if err != nil {
		return payload.WrapBadRequestErr(err)
	}

	return h.Storage().Delete(&models.Account{ID: &ID})
}
