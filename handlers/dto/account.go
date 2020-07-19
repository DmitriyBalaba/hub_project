package dto

import (
	"hub_project/models"

	"bitbucket.org/jazzserve/webapi/web/http/payload"
)

type AccountPost struct {
	Name           string        `json:"name"`
	Email          string        `json:"email" validate:"email,required"`
	IsAdmin        bool          `json:"is_admin"`
	IsStoreManager bool          `json:"is_store_manager"`
	Store          *models.Store `json:"store,omitempty"`
}

type AccountPut struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	IsAdmin        bool   `json:"is_admin"`
	IsStoreManager bool   `json:"is_store_manager"`
}

type AccountGet struct {
	payload.OffsetLimit
}
