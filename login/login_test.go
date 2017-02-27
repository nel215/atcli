package login

import (
	"net/http"
	"testing"
)

func postForTesing(name string, password string) (*http.Response, error) {
	header := http.Header{}
	cookies := []http.Cookie{
		http.Cookie{Name: "_session", Value: "_session"},
		http.Cookie{Name: "_issue_time", Value: "_issue_time"},
		http.Cookie{Name: "_kick_id", Value: "_kick_id"},
		http.Cookie{Name: "_user_id", Value: "_user_id"},
	}

	for _, c := range cookies {
		header.Add("Set-Cookie", c.String())
	}
	resp := &http.Response{
		StatusCode: 200,
		Header:     header,
	}

	return resp, nil
}

func TestSubmit(t *testing.T) {
	l := &Login{postForTesing}
	s, err := l.Submit("name", "password")
	if err != nil {
		t.Fatalf("submit failed")
	}

	if s.Session != "_session" {
		t.Fatalf("s.Session expected _session . got %s", s.Session)
	}
	if s.IssueTime != "_issue_time" {
		t.Fatalf("s.IssueTim expected _issue_time. got %s", s.IssueTime)
	}
	if s.KickId != "_kick_id" {
		t.Fatalf("s.KickId expected _kick_id. got %s", s.KickId)
	}
	if s.UserID != "_user_id" {
		t.Fatalf("s.UserID expected _user_id. got %s", s.UserID)
	}
}
