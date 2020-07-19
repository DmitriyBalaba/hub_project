package dto

import "hub_project/models"

type AccountStorePost struct {
	Account *models.Account `json:"account,omitempty"  validate:"required"`
	Store   *models.Store   `json:"store,omitempty" validate:"required"`
}
