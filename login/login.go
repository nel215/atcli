package login

import (
	"errors"
	"net/http"
	"net/url"
)

type Login struct {
	post func(name string, password string) (*http.Response, error)
}

func New() *Login {
	return &Login{
		post: post,
	}
}

func post(name string, password string) (*http.Response, error) {
	data := url.Values{}
	data.Set("name", name)
	data.Set("password", password)
	return http.DefaultClient.PostForm("https://practice.contest.atcoder.jp/login", data)
}

type Session struct {
	Session   string
	IssueTime string
	KickId    string
	UserID    string
}

func NewSessionFromCookies(cookies []*http.Cookie) *Session {
	s := &Session{}
	for _, c := range cookies {
		switch c.Name {
		case "_session":
			s.Session = c.Value
		case "_issue_time":
			s.IssueTime = c.Value
		case "_kick_id":
			s.KickId = c.Value
		case "_user_id":
			s.UserID = c.Value
		}
	}
	return s
}

func (l *Login) Submit(name string, password string) (*Session, error) {
	resp, err := l.post(name, password)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 302 {
		return nil, errors.New("login failed")
	}

	s := NewSessionFromCookies(resp.Cookies())

	return s, nil
}
