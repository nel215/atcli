package api

import (
	"errors"
	"fmt"
	"github.com/nel215/atcli/session"
	"github.com/nel215/atcli/store"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

type Submit struct {
	problemId      int64
	languageId     int64
	sourceCodePath string
	sessionStore   interface {
		Load() (*session.Session, error)
	}
	configStore interface {
		Load() (*store.Config, error)
	}
}

func NewSubmit(problemId int64, languageId int64, sourceCodePath string) (*Submit, error) {
	if problemId == 0 {
		return nil, errors.New("problemId is required")
	}
	if languageId == 0 {
		return nil, errors.New("languageId is required")
	}
	if sourceCodePath == "" {
		return nil, errors.New("sourceCodePath is required")
	}
	return &Submit{
		problemId:      problemId,
		languageId:     languageId,
		sourceCodePath: sourceCodePath,
		sessionStore:   store.NewSessionStore(),
		configStore:    store.NewConfigStore(),
	}, nil
}

func (s *Submit) createForm(csrfToken string) (url.Values, error) {
	f, err := os.Open(s.sourceCodePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sourceCode, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("__session", csrfToken)
	data.Set("task_id", fmt.Sprintf("%d", s.problemId))
	data.Set(fmt.Sprintf("language_id_%d", s.problemId), fmt.Sprintf("%d", s.languageId))
	data.Set("source_code", string(sourceCode))
	return data, nil
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

func (s *Submit) Execute() error {

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

	data, err := s.createForm(csrfToken)
	if err != nil {
		return err
	}
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
