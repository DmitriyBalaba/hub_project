package dto

import (
	"hub_project/models"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
)

type StorePost struct {
	models.Store
	AccountList []*models.Account `json:"account_list,omitempty"`
}

type StorePut StorePost

type AccountStoreQuery struct {
	AccountID *int64 `schema:"account_id,required"`
	payload.OffsetLimit
}
