package memory

import (
	"github.com/gorilla/securecookie"
	gorilla "github.com/gorilla/sessions"
	"net/http"
)

func (s *SessionStorage) Get(r *http.Request, name string) (*gorilla.Session, error) {
	return gorilla.GetRegistry(r).Get(s, name)
}

func (s *SessionStorage) New(r *http.Request, name string) (*gorilla.Session, error) {
	session := gorilla.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	cookie, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}

	if err = securecookie.DecodeMulti(name, cookie.Value, &session.ID, s.Codecs...); err != nil {
		return session, nil
	}

	memValue, ok := s.GetSessionValue(session.ID)
	if !ok {
		return session, nil
	}

	session.Values = copySessionValues(memValue)
	session.IsNew = false
	return session, nil
}

func (s *SessionStorage) Save(r *http.Request, w http.ResponseWriter, session *gorilla.Session) (err error) {
	var cookie string
	if session.Options.MaxAge < 0 {
		s.deleteSessionValue(session)
	} else {
		if session.ID == "" {
			session.ID = s.genSessionID(session)
		}
		if cookie, err = securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...); err != nil {
			return
		}
		s.saveSessionValue(session)
	}
	http.SetCookie(w, gorilla.NewCookie(session.Name(), cookie, session.Options))
	return nil
}
