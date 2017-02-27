package store

import (
	"encoding/json"
	"github.com/nel215/atcmd/login"
	"os"
)

type SessionStore struct{}

func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

func (ss *SessionStore) Save(sess *login.Session) error {
	jb, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	f, err := os.Create("./.session.json")
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(jb)
	return nil
}
