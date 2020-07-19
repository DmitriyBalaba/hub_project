package sendgrid

import "github.com/pkg/errors"

type template struct {
	ID string `validate:"required" yaml:"id"`
}

type templates map[string]template

type Config struct {
	ApiKey    string    `validate:"required" yaml:"api-key"`
	From      string    `validate:"required" yaml:"from"`
	FromName  string    `validate:"required" yaml:"from-name"`
	Templates templates `validate:"dive" yaml:"templates"`
}

var requiredTemplates = make(map[string]struct{})

func AddRequiredTemplates(keys ...string) {
	for i := range keys {
		requiredTemplates[keys[i]] = struct{}{}
	}
}

func checkRequired(templates templates) error {
	for k := range requiredTemplates {
		if _, ok := templates[k]; !ok {
			return errors.Errorf("required %s send grid template not found", k)
		}
	}
	return nil
}
