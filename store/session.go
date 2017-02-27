package store

import (
	"encoding/json"
	"github.com/nel215/atcmd/login"
	"io/ioutil"
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

func (ss *SessionStore) Load() (*login.Session, error) {
	f, err := os.Open("./.session.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	sess := &login.Session{}
	err = json.Unmarshal(bs, sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
