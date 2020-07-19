package store

import (
	"crypto/sha512"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var NilConfigErr = errors.New("can't create new session store: nil config")

func New(c *Config) (store *SessionStore, err error) {
	if c == nil {
		err = NilConfigErr
		return
	}

	if *c.MaxAge <= time.Minute {
		log.Warn().Msgf("Session Max Age is too low: %v! Is it a typo?", c.MaxAge)
	}

	store = &SessionStore{
		pool:   newPool(c),
		codecs: securecookie.CodecsFromPairs([]byte(c.Key)),
		options: &sessions.Options{
			Path: defaultCookiePath,
		},
		GetSessionID: GetSessionIDWithUserID,
		Interpreter:  &jsonInterpreter{},
	}

	store.MaxAge(int(c.MaxAge.Seconds()))
	return
}

func newPool(c *Config) *redis.Pool {
	dialFunc := func() (redis.Conn, error) {
		conn, err := redis.Dial(network, c.Address)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if _, err := conn.Do(Select, c.DB); err != nil {
			conn.Close()
			return nil, errors.WithStack(err)
		}
		return conn, nil
	}

	return &redis.Pool{
		MaxIdle:     defaultMaxIdle,
		IdleTimeout: defaultIdleTimeOut,
		Dial:        dialFunc,
	}
}

type SessionStore struct {
	pool         *redis.Pool
	codecs       []securecookie.Codec
	options      *sessions.Options
	GetSessionID GetSessIDFunc
	Interpreter  Interpreter
}

type GetSessIDFunc func(*http.Request) (string, error)

func (s *SessionStore) MaxAge(maxAge int) {
	s.options.MaxAge = maxAge
	for _, s := range s.codecs {
		if cookie, ok := s.(*securecookie.SecureCookie); ok {
			cookie.MaxAge(maxAge)
			cookie.SetSerializer(securecookie.JSONEncoder{})
			cookie.HashFunc(sha512.New512_256)
		}
	}
}

func (s *SessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *SessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := s.newSession(name)
	if c, errCookie := r.Cookie(name); errCookie == nil {
		errDecode := securecookie.DecodeMulti(name, c.Value, &session.ID, s.codecs...)
		if errDecode == nil {
			isDataAvailable, errRedis := s.load(r, session)
			session.IsNew = errRedis != nil && !isDataAvailable
			return session, errRedis
		}
	}
	return session, nil
}

func (s *SessionStore) retrieveSession(sessionId string) (interface{}, error) {
	conn := s.pool.Get()
	defer conn.Close()

	err := conn.Send(Expire, sessionId, s.options.MaxAge)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: retrieve only certain fields based on obj type structure
	reply, err := conn.Do(HGetAll, sessionId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reply, nil
}

func (s *SessionStore) load(r *http.Request, session *sessions.Session) (bool, error) {
	redisValues, err := redis.StringMap(s.retrieveSession(session.ID))
	if err != nil {
		return false, errors.WithStack(err)
	}

	if len(redisValues) == 0 {
		return false, nil
	}

	objType := s.Interpreter.GetObjType(r)
	obj, err := s.Interpreter.ConvertObjFromRedis(redisValues, objType)
	if err != nil {
		return true, err
	}

	err = s.Interpreter.PutObj(session, obj)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (s *SessionStore) newSession(name string) (session *sessions.Session) {
	session = sessions.NewSession(s, name)
	opts := *s.options
	session.Options = &opts
	session.IsNew = true
	return
}

func (s *SessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		if err := s.save(r, session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
		if err != nil {
			return errors.WithStack(err)
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *SessionStore) delete(session *sessions.Session) error {
	_, err := s.doSingle(Del, session.ID)
	return errors.WithStack(err)
}

func (s *SessionStore) save(r *http.Request, session *sessions.Session) (err error) {
	if session.IsNew || session.ID == "" {
		session.ID, err = s.GetSessionID(r)
		if err != nil {
			return
		}
	}

	redisValues, err := s.Interpreter.ConvertObjToRedis(s.Interpreter.GetObj(session))
	if err != nil {
		return err
	}

	conn := s.pool.Get()
	defer conn.Close()

	err = conn.Send(HMSet, redis.Args{}.Add(session.ID).AddFlat(redisValues)...)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = conn.Do(Expire, session.ID, session.Options.MaxAge)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *SessionStore) DoesSessionExist(sessionID string) (bool, error) {
	ok, err := redis.Bool(s.doSingle(Exists, sessionID))
	return ok, errors.WithStack(err)
}

func (s *SessionStore) UpdateSession(sessionID string, v interface{}) error {
	redisValues, err := s.Interpreter.ConvertObjToRedis(v)
	if err != nil {
		return err
	}

	_, err = s.doSingle(HMSet, redis.Args{}.Add(sessionID).AddFlat(redisValues)...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *SessionStore) UpdateAllSessions(pattern string, v interface{}) error {
	redisValues, err := s.Interpreter.ConvertObjToRedis(v)
	if err != nil {
		return err
	}

	getArgs := func(key string) (redis.Args, error) {
		return redis.Args{}.Add(key).AddFlat(redisValues), nil
	}

	return s.updateAllSessions(pattern, getArgs)
}

func (s *SessionStore) newCache(getter ObjGetter) *valuesCache {
	return &valuesCache{
		values:      make(map[interface{}]map[string]string),
		getter:      getter,
		interpreter: s.Interpreter,
	}
}

func (s *SessionStore) UpdateAllSessionsDynamic(pattern string, getter ObjGetter) error {
	return s.updateAllSessions(pattern, s.newCache(getter).getArgs)
}

func (s *SessionStore) doSingle(commandName string, args ...interface{}) (interface{}, error) {
	conn := s.pool.Get()
	defer conn.Close()
	reply, err := conn.Do(commandName, args...)
	return reply, errors.WithStack(err)
}
