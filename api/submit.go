package api

import (
	"errors"
	"fmt"
	"github.com/nel215/atcli/session"
	"github.com/nel215/atcli/store"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Submit struct {
	problemId    int64
	sessionStore interface {
		Load() (*session.Session, error)
	}
	configStore interface {
		Load() (*store.Config, error)
	}
}

func NewSubmit(problemId int64) (*Submit, error) {
	if problemId == 0 {
		return nil, errors.New("problemId is required")
	}
	return &Submit{
		problemId:    problemId,
		sessionStore: store.NewSessionStore(),
		configStore:  store.NewConfigStore(),
	}, nil
}

func createForm(problemId int64, languageId int64, sourceCode []byte, csrfToken string) url.Values {
	data := url.Values{}
	data.Set("__session", csrfToken)
	data.Set("task_id", fmt.Sprintf("%d", problemId))
	data.Set(fmt.Sprintf("language_id_%d", problemId), fmt.Sprintf("%d", languageId))
	data.Set("source_code", string(sourceCode))
	return data
}

func extractCSRFToken(sess *session.Session, contestUrl string) (string, error) {
	resp, err := http.Get(contestUrl + "/submit")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return "", nil
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data != "input" {
				continue
			}
			isCSRFToken := false
			for _, a := range t.Attr {
				if a.Key == "name" && a.Val == "__session" {
					isCSRFToken = true
				}
			}
			if !isCSRFToken {
				continue
			}
			for _, a := range t.Attr {
				if a.Key == "value" {
					return a.Val, nil
				}
			}
		}
	}
	return "", nil
}

func (s *Submit) Execute(languageId int64, sourceCode []byte) error {
	sess, err := s.sessionStore.Load()
	if err != nil {
		return err
	}
	config, err := s.configStore.Load()
	if err != nil {
		return err
	}
	contestUrl := config.ContestUrl
	jar, err := cookiejar.New(nil)
	cookies := sess.GetSessionCookies()
	u, err := url.Parse(contestUrl)
	if err != nil {
		return err
	}
	jar.SetCookies(u, cookies)
	http.DefaultClient.Jar = jar
	csrfToken, err := extractCSRFToken(sess, contestUrl)
	if err != nil {
		return err
	}
	data := createForm(s.problemId, languageId, sourceCode, csrfToken)
	submitUrl := fmt.Sprintf("%s/submit?task_id=%d", contestUrl, s.problemId)
	resp, err := http.PostForm(submitUrl, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		m := fmt.Sprintf("expected StatusCode is 200. got %d", resp.StatusCode)
		return errors.New(m)
	}
	return nil
}
