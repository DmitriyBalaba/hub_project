package env

import (
	"hub_project/models"

	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/rs/zerolog/log"
)

func (e *env) SendNotificationToNewAccount(acc *models.Account) (err error) {

	data := &struct {
		Email    *string
		Name     *string
		Password *string
		Link     string
	}{Email: acc.Email, Name: acc.Name, Password: acc.Password, Link: e.GetRootURL()}

	tmpl, err := e.Email().HTMLTemplateToString("templateNotificationCreateAccount.html", data)
	if err != nil {
		return err
	}

	err = e.Email().
		SendHTMLMessage([]string{val.Str(acc.Email)}, "[Savvie.io] Account Creation", tmpl, nil)
	if err != nil {
		log.Error().Msgf("notify about account locked is failed [%s]", err.Error())
	}

	return nil
}
