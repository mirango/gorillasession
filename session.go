package gorillasession

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mirango/framework"
)

type Session struct {
	store   *Store // maybe multiple stores
	session *sessions.Session
	changed bool
}

func (s *Session) Name() string {
	return s.session.Name()
}

func (s *Session) ID() string {
	return s.session.ID
}

func (s *Session) Values() framework.Values {
	return framework.Values(s.session.Values)
}

func (s *Session) Changed() bool {
	return s.changed
}

func (s *Session) Store() framework.SessionStore {
	return s.store
}

func (s *Session) Get(key interface{}) framework.Value {
	return framework.NewValue(s.session.Values[key])
}

func (s *Session) GetOr(key interface{}, or interface{}) framework.Value {
	v, ok := s.session.Values[key]
	if !ok {
		return framework.NewValue(or)
	}
	return framework.NewValue(v)
}

func (s *Session) Set(key interface{}, val interface{}) {
	s.changed = true
	s.session.Values[key] = val
}

func (s *Session) Unset(key interface{}) bool {
	_, isset := s.session.Values[key]
	if isset {
		delete(s.session.Values, key)
		s.changed = true
	}
	return isset
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	if s.changed {
		return s.session.Save(r, w)
	}
	return nil
}

func (s *Session) Delete(r *http.Request, w http.ResponseWriter) error {
	return nil
}

func (s *Session) Flashes(key ...string) []framework.Value {
	flashes := s.session.Flashes(key...)
	if len(flashes) > 0 {
		s.changed = true
		return nil
	}
	return nil
}

func (s *Session) AddFlash(val interface{}, key ...string) {
	s.changed = true
	s.session.AddFlash(val, key...)
}
