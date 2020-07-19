package env

import (
	"bitbucket.org/jazzserve/webapi/databases/gorm/postgres"
	"bitbucket.org/jazzserve/webapi/email/sendgrid"
	"bitbucket.org/jazzserve/webapi/email/smtp"
	"bitbucket.org/jazzserve/webapi/resources/filesys"
	"bitbucket.org/jazzserve/webapi/security/jwt"
	"bitbucket.org/jazzserve/webapi/sessions/memory"
	"bitbucket.org/jazzserve/webapi/sessions/redis/store"
	"bitbucket.org/jazzserve/webapi/utils"
	"bitbucket.org/jazzserve/webapi/web/http/cors"
	"bitbucket.org/jazzserve/webapi/web/http/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

type Config struct {
	AppName      *string `yaml:"app-name"`
	SessionStore struct {
		Redis  *store.Config  `yaml:"redis"`
		Memory *memory.Config `yaml:"memory"`
	} `yaml:"session-store"`
	Database struct {
		Postgres *postgres.Config `yaml:"postgres"`
	} `yaml:"database"`
	Resources struct {
		FileSystem *filesys.Config `yaml:"file-system"`
	} `yaml:"resources"`
	Email struct {
		Smtp     *smtp.Config     `yaml:"smtp"`
		SendGrid *sendgrid.Config `yaml:"send-grid"`
	} `yaml:"email"`
	Security struct {
		JWT *jwt.Config `yaml:"jwt"`
	} `yaml:"security"`
	Server *server.Config `yaml:"server"`
	CORS   *cors.Config   `yaml:"cors"`
}

func ReadValidYamlFile(dst interface{}, filePath string) (err error) {
	cType := reflect.TypeOf(dst)
	if cType.Kind() != reflect.Ptr || cType.Elem().Kind() != reflect.Struct {
		return errors.New("expected config as struct ptr")
	}

	log.Info().Msgf("Reading configuration from: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer utils.CheckCloseError(file, &err)

	if err = yaml.NewDecoder(file).Decode(dst); err != nil {
		err = errors.Wrap(err, "can't decode config file")
		return
	}

	if err = validator.New().Struct(dst); err != nil {
		err = errors.Wrap(err, "can't validate config file")
		return
	}

	return
}
