package gorillasession

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mirango/framework"
)

type Store struct {
	store sessions.Store
	names []string
}

func New(store sessions.Store, names ...string) *Store {
	return &Store{store: store, names: names}
}

func (s *Store) Names() []string {
	return s.names
}

func (s *Store) Get(r *http.Request, name string) (framework.Session, error) {
	if !s.validName(name) {
		return nil, nil
	}
	se, err := s.store.Get(r, name)
	if err != nil {
		return nil, err
	}
	return &Session{s, se, false}, nil
}

func (s *Store) GetAll(r *http.Request) (framework.Sessions, error) {
	return s.GetMany(r, s.names...)
}

func (s *Store) GetMany(r *http.Request, names ...string) (framework.Sessions, error) {
	var sess framework.Sessions
	for _, name := range s.names {
		se, err := s.Get(r, name)
		if err != nil {
			return nil, err
		}
		sess.Append(se)
	}
	return sess, nil
}

func (s *Store) New(r *http.Request, name string) (framework.Session, error) {
	if !s.validName(name) {
		return nil, nil
	}
	se, err := s.store.New(r, name)
	if err != nil {
		return nil, err
	}
	return &Session{s, se, false}, nil
}

func (s *Store) validName(name string) bool {
	for _, sn := range s.names {
		if name == sn {
			return true
		}
	}
	return false
}

func (s *Store) validStore(store framework.SessionStore) bool {
	// for _, sn := range s.sessions {
	// 	if name == sn {
	// 		return true
	// 	}
	// }
	if store == s {
		return true
	}
	return false
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, se framework.Session) error {
	if !s.validName(se.Name()) || !s.validStore(se.Store()) {
		return nil
	}

	sess, ok := se.(*Session)
	if !ok {
		return nil
	}

	err := s.store.Save(r, w, sess.session)
	if err != nil {
		return nil
	}
	return nil
}
