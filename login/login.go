package login

import (
	"errors"
	"fmt"
	"github.com/nel215/atcli/session"
	"github.com/nel215/atcli/store"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Login struct {
	sessionStore interface {
		Save(*session.Session) error
	}
	post func(name string, password string) (*http.Response, error)
}

func New() *Login {
	return &Login{
		sessionStore: store.NewSessionStore(),
		post:         post,
	}
}

func post(name string, password string) (*http.Response, error) {
	data := url.Values{}
	data.Set("name", name)
	data.Set("password", password)
	return http.DefaultClient.PostForm("https://practice.contest.atcoder.jp/login", data)
}

func (l *Login) Submit(user string, password string) error {
	if user == "" {
		return errors.New("user is required")
	}
	if password == "" {
		return errors.New("password is required")
	}

	log.Printf("try logging in by %s...\n", user)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	http.DefaultClient.Jar = jar
	resp, err := l.post(user, password)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("expected StatusCode is 200. got %d", resp.StatusCode)
		return errors.New(m)
	}

	u, err := url.Parse("https://practice.contest.atcoder.jp")
	if err != nil {
		return err
	}
	sess := session.NewSessionFromCookies(jar.Cookies(u))

	err = l.sessionStore.Save(sess)
	if err != nil {
		return err
	}

	return nil
}
