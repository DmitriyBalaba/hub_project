package memory

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/securecookie"
	gorilla "github.com/gorilla/sessions"
	"regexp"
	"strings"
	"sync"
)

type (
	SessionID          = string
	SessionValue       = map[interface{}]interface{}
	sessionValues      = map[SessionID]SessionValue
	sessionIDGenerator = func(session *gorilla.Session) string
	SessionStorage     struct {
		Codecs       []securecookie.Codec
		Options      *gorilla.Options
		values       sessionValues
		rwMutex      sync.RWMutex
		genSessionID sessionIDGenerator
	}
)

func NewSessionStorage(config Config) *SessionStorage {
	s := &SessionStorage{
		Codecs:       securecookie.CodecsFromPairs([]byte(config.Key)),
		Options:      &gorilla.Options{Path: "/"},
		values:       make(sessionValues),
		genSessionID: defaultSessionIDGenerator,
	}
	s.setMaxAge(config.MaxAge)
	return s
}

func (s *SessionStorage) SetSessionIDGenerator(f sessionIDGenerator) {
	s.genSessionID = f
}

func (s *SessionStorage) setMaxAge(age int) {
	s.Options.MaxAge = age
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (s *SessionStorage) saveSessionValue(session *gorilla.Session) {
	if session.Values == nil {
		panic("cannot save session with empty Values")
	}
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	s.values[session.ID] = copySessionValues(session.Values)
}

// deletes all values in the current session's values storage
// and the session's value that was stored in the memory storage
func (s *SessionStorage) deleteSessionValue(session *gorilla.Session) {
	for k := range session.Values {
		delete(session.Values, k)
	}
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	if _, ok := s.values[session.ID]; ok {
		delete(s.values, session.ID)
	}
}

// retrieves the in-memory value of the given session id
func (s *SessionStorage) GetSessionValue(id SessionID) (SessionValue, bool) {
	if s.values == nil {
		return nil, false
	}
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	value, ok := s.values[id]
	return value, ok
}

func copySessionValues(v SessionValue) (vCopy SessionValue) {
	var b = &bytes.Buffer{}

	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		panic(fmt.Errorf("failed to copy session value: %s", err.Error()))
	}

	err = gob.NewDecoder(b).Decode(&vCopy)
	if err != nil {
		panic(fmt.Errorf("failed to copy session value: %s", err.Error()))
	}

	return
}

func defaultSessionIDGenerator(session *gorilla.Session) string {
	return strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
}

func (s *SessionStorage) DeleteAllSessionsWithPattern(pattern *regexp.Regexp) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	for key := range s.values {
		if pattern.MatchString(key) {
			delete(s.values, key)
		}
	}
}

func (s *SessionStorage) UpdateAllSessionsWithPattern(p *regexp.Regexp, v SessionValue) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	for key := range s.values {
		if p.MatchString(key) {
			s.values[key] = v
		}
	}
}
