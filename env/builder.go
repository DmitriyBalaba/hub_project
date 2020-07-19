package env

import (
	"context"
	"errors"
	"hub_project/models"

	"bitbucket.org/jazzserve/webapi/databases/gorm/postgres"
	"bitbucket.org/jazzserve/webapi/email/smtp"
	lib_env "bitbucket.org/jazzserve/webapi/env"
	"bitbucket.org/jazzserve/webapi/resources/filesys"
	"bitbucket.org/jazzserve/webapi/sessions/redis/store"
	"bitbucket.org/jazzserve/webapi/web/http/cors"
	"bitbucket.org/jazzserve/webapi/web/http/routing"
	"bitbucket.org/jazzserve/webapi/web/http/server"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/justinas/alice"
)

type Builder struct {
	*env
	middleware []alice.Constructor
	debugMode  bool
	err        error
}

func NewBuilder(debugMode bool) *Builder {
	return &Builder{
		env:       &env{},
		debugMode: debugMode,
	}
}

func (b *Builder) ReadConfigFile(filePath string) *Builder {
	if b.hasErr() {
		return b
	}

	if b.config != nil {
		b.err = errors.New("config already set")
		return b
	}

	b.config = new(lib_env.Config)
	b.err = lib_env.ReadValidYamlFile(b.config, filePath)
	return b
}

func getPostgresDB(config *postgres.Config, debugMode bool) (db *gorm.DB, err error) {
	db, err = postgres.NewDB(config, debugMode)
	if err != nil {
		return
	}

	db.SingularTable(true)
	return
}

func (b *Builder) SetupPostgresDB() *Builder {
	if !b.validConfig() {
		return b
	}

	db, err := getPostgresDB(b.config.Database.Postgres, b.debugMode)
	if err != nil {
		b.err = err
		return b
	}

	b.storage = models.NewStorage(db)
	return b
}

func (b *Builder) SetupSessionStore() *Builder {
	if !b.validConfig() {
		return b
	}

	b.sessionStore, b.err = store.New(b.config.SessionStore.Redis)
	return b
}

func (b *Builder) SetupEmail() *Builder {
	if !b.validConfig() {
		return b
	}

	b.email = smtp.NewEmail(b.config.Email.Smtp)
	if b.fileManager != nil {
		b.email.Templates, b.err = filesys.ParseTemplates(b.fileManager.GetDirPath(TemplatesDir))
	} else {
		panic("resource manager not built yet")
	}
	return b
}

func (b *Builder) SetupFileManager(directories ...string) *Builder {
	if !b.validConfig() {
		return b
	}

	filesys.SetRequiredDirs(directories...)
	b.fileManager, b.err = filesys.NewManager(b.config.Resources.FileSystem)
	if b.err != nil {
		return b
	}

	b.err = b.fileManager.CreateDirectories()
	return b
}

func (b *Builder) ConfigureCors() *Builder {
	if !b.validConfig() {
		return b
	}

	corsMiddleware := cors.
		New(b.config.CORS, AllowedHeaders(), AllowedMethods()).
		Middleware()

	return b.AddGlobalMiddleware(corsMiddleware)
}

func (b *Builder) AddGlobalMiddleware(middle ...alice.Constructor) *Builder {
	if b.hasErr() {
		return b
	}
	b.middleware = append(b.middleware, middle...)
	return b
}

type RouterConfig func(r *mux.Router) *mux.Router

func (b *Builder) BuildWebServer(routerConfigs ...RouterConfig) (*server.Server, error) {
	environment, err := b.Build()
	if err != nil {
		return nil, err
	}

	server := server.New(b.config.Server)

	router := routing.NewRouter(b.config.Server.GetPathPrefix())
	for _, configure := range routerConfigs {
		router = configure(router)
	}

	routing.AddRoutesTo(router, Add(context.Background(), environment))

	handler := alice.New(b.middleware...).Then(router)
	server.SetHandler(handler)

	return server, nil
}

var BuiltAlreadyErr = errors.New("already built")

func (b *Builder) Build() (*env, error) {
	if b.hasErr() {
		return nil, b.err
	}

	b.err = nil
	return b.env, nil
}

func (b *Builder) validConfig() bool {
	if b.hasErr() {
		return false
	}
	if b.config == nil {
		b.err = NilConfigErr
		return false
	}
	return true
}

func (b *Builder) hasErr() bool {
	if b.err != nil {
		return true
	}
	return false
}
