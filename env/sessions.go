package env

import (
	"hub_project/models"
	"net/http"
	"strconv"

	"bitbucket.org/jazzserve/webapi/sessions/redis/store"
	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

func (e *env) GetSession(r *http.Request) (*sessions.Session, *http.Request, error) {
	if e.sessionStore == nil {
		panic("can't get session with nil session store")
	}

	if r == nil {
		return nil, nil, errors.New("can't proceed with zero input")
	}

	r, err := e.sessionStore.Interpreter.PutObjType(r, &models.Account{})
	if err != nil {
		return nil, nil, err
	}

	session, err := e.sessionStore.Get(r, SessionIDKey)
	if err != nil {
		return nil, nil, err
	}

	return session, r, nil
}

func (e *env) isSessionNil() error {
	if e.sessionStore == nil || e.storage == nil {
		return errors.New("sessions are nil")
	}
	return nil
}

func (e *env) GetFromSession(session *sessions.Session) (*models.Account, bool) {
	if e.sessionStore == nil {
		panic("can't get from session with nil session store")
	}

	v := e.sessionStore.Interpreter.GetObj(session)
	if v == nil {
		return nil, false
	}

	return v.(*models.Account), true
}

func (e *env) AddToSession(a *models.Account, session *sessions.Session) error {
	if e.sessionStore == nil {
		panic("can't add to session with nil session store")
	}

	return e.sessionStore.Interpreter.PutObj(session, a)
}

func (e *env) SaveSessionAsNew(r *http.Request, w http.ResponseWriter, s *sessions.Session, userID *int64) error {
	if e.sessionStore == nil {
		panic("can't save session with nil session store")
	}

	r, err := store.PutUserID(r, strconv.Itoa(int(val.Int64(userID))))
	if err != nil {
		return err
	}

	return errors.WithStack(s.Save(r, w))
}

func (e *env) DeleteSession(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	if e.sessionStore == nil {
		panic("can't delete session with nil session store")
	}

	s.Options.MaxAge = -1
	if err := s.Save(r, w); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
