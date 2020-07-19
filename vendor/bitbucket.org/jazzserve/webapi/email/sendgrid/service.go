package sendgrid

import (
	"bitbucket.org/jazzserve/webapi/email"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

func New(c *Config) *Service {
	if c == nil {
		panic("can't create sendgrid service with nil config")
	}

	if err := checkRequired(c.Templates); err != nil {
		panic(err)
	}

	return &Service{
		config: *c,
		client: sendgrid.NewSendClient(c.ApiKey),
	}
}

type Service struct {
	config Config
	client *sendgrid.Client
}

func (s *Service) NewMail() email.Email {
	return &Mail{
		service: s,
		mail: mail.
			NewV3Mail().
			SetFrom(mail.NewEmail(s.config.FromName, s.config.From)),
		personalization: mail.NewPersonalization(),
	}
}

type Mail struct {
	service         *Service
	mail            *mail.SGMailV3
	personalization *mail.Personalization
}

func (m *Mail) SetTemplate(template string) email.Email {
	templateID, ok := m.service.config.Templates[template]
	if !ok {
		panic(fmt.Sprintf("sendgrid template '%s' not found", template))
	}

	m.mail = m.mail.SetTemplateID(templateID.ID)
	return m
}

func (m *Mail) SetTo(name string, address string) email.Email {
	m.personalization.AddTos(mail.NewEmail(name, address))
	return m
}

func (m *Mail) SetDynamicData(key string, value string) email.Email {
	m.personalization.SetDynamicTemplateData(key, value)
	return m
}

func (m *Mail) Send() (err error) {
	sgMail := m.mail.AddPersonalizations(m.personalization)
	r, err := m.service.client.Send(sgMail)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusAccepted {
		return errors.New(r.Body)
	}

	return
}
