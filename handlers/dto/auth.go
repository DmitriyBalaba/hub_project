package dto

type AccountLogin struct {
	Email    string `json:"email,omitempty" validate:"email,required"`
	Password string `json:"password,omitempty" validate:"required"`
}
