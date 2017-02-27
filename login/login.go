package login

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
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

func (l *Login) Submit(name string, password string) (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	http.DefaultClient.Jar = jar
	resp, err := l.post(name, password)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("expected StatusCode is 200. got %d", resp.StatusCode)
		return nil, errors.New(m)
	}

	u, err := url.Parse("https://practice.contest.atcoder.jp")
	if err != nil {
		return nil, err
	}
	s := NewSessionFromCookies(jar.Cookies(u))

	return s, nil
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

func (s *Session) AddSessionCookies(req *http.Request) {
	req.AddCookie(&http.Cookie{Name: "_session", Value: s.Session})
	req.AddCookie(&http.Cookie{Name: "_issue_time", Value: s.IssueTime})
	req.AddCookie(&http.Cookie{Name: "_kick_id", Value: s.KickId})
	req.AddCookie(&http.Cookie{Name: "_user_id", Value: s.UserID})
}
