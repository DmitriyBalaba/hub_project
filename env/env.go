package env

import (
	"hub_project/models"
	"net/http"

	"bitbucket.org/jazzserve/webapi/email/smtp"
	. "bitbucket.org/jazzserve/webapi/env"
	"bitbucket.org/jazzserve/webapi/resources/filesys"
	"bitbucket.org/jazzserve/webapi/sessions/redis/store"
	"github.com/gorilla/sessions"
)

const (
	MigrationDir = "migration"
	TemplatesDir = "templates"
)

type Environment interface {

	// env.go
	Storage() models.Storage
	Email() *smtp.Email
	ResourceManager() *filesys.Manager
	GetRootURL() string

	//sessions.go
	GetSession(r *http.Request) (*sessions.Session, *http.Request, error)
	GetFromSession(session *sessions.Session) (*models.Account, bool)
	AddToSession(a *models.Account, session *sessions.Session) error
	SaveSessionAsNew(r *http.Request, w http.ResponseWriter, s *sessions.Session, userID *int64) error
	DeleteSession(r *http.Request, w http.ResponseWriter, s *sessions.Session) error

	// auth.go
	AuthAdmin(next http.Handler) http.Handler
	Auth(next http.Handler) http.Handler

	// migration.go
	ReadMigrationScripts() (ms []*models.Migration, err error)

	//email.go
	SendNotificationToNewAccount(acc *models.Account) (err error)
}

type env struct {
	storage      models.Storage
	config       *Config
	fileManager  *filesys.Manager
	email        *smtp.Email
	sessionStore *store.SessionStore
}

func (e *env) GetRootURL() string {
	if e.config == nil || len(e.config.CORS.Origins) == 0 {
		panic("nil CORS origins")
	}
	return e.config.CORS.Origins[0]
}

func (e *env) ResourceManager() *filesys.Manager {
	return e.fileManager
}

// TODO: remove panic from getters

func (e *env) Storage() models.Storage {
	if e.storage == nil {
		panic("nil storage")
	}
	return e.storage
}

func (e *env) Email() *smtp.Email {
	if e.email == nil {
		panic("nil email")
	}
	return e.email
}

func (e *env) CopyWithNewStorage(newStorage models.Storage) *env {
	newEnv := *e
	newEnv.storage = newStorage
	return &newEnv
}
