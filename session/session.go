package session

import (
	"net/http"
)

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

func (s *Session) GetSessionCookies() []*http.Cookie {
	return []*http.Cookie{
		&http.Cookie{Name: "_session", Value: s.Session},
		&http.Cookie{Name: "_issue_time", Value: s.IssueTime},
		&http.Cookie{Name: "_kick_id", Value: s.KickId},
		&http.Cookie{Name: "_user_id", Value: s.UserID},
	}
}
