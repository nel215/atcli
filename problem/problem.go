package problem

import (
	"fmt"
	"github.com/nel215/atcli/session"
	"github.com/nel215/atcli/store"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Problem struct {
	sessionStore interface {
		Load() (*session.Session, error)
	}
	configStore interface {
		Load() (*store.Config, error)
	}
}

func New() *Problem {
	return &Problem{
		store.NewSessionStore(),
		store.NewConfigStore(),
	}
}

func (p *Problem) Execute() error {
	config, err := p.configStore.Load()
	if err != nil {
		return err
	}
	contestUrl := config.ContestUrl
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	u, err := url.Parse(contestUrl)
	if err != nil {
		return err
	}
	sess, err := p.sessionStore.Load()
	if err != nil {
		return err
	}
	cookies := sess.GetSessionCookies()
	jar.SetCookies(u, cookies)
	http.DefaultClient.Jar = jar

	resp, err := http.Get(fmt.Sprintf("%s/submit", contestUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	InTaskSelector := false
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "select" {
				InTaskSelector = checkTaskSelector(t)
				if InTaskSelector {
					fmt.Printf("# Problems\n\n")
				}
			}

			if !InTaskSelector {
				continue
			}

			if t.Data == "option" {
				id := extractTaskId(t)
				z.Next()
				t = z.Token()
				fmt.Printf("- %s: %s\n", id, t.Data)
			}
		case tt == html.EndTagToken:
			t := z.Token()

			if t.Data == "select" {
				InTaskSelector = false
			}
		}
	}
}

func checkTaskSelector(t html.Token) bool {
	for _, a := range t.Attr {
		if a.Key == "name" && a.Val == "task_id" {
			return true
		}
	}
	return false
}

func extractTaskId(t html.Token) string {
	for _, a := range t.Attr {
		if a.Key == "value" {
			return a.Val
		}
	}
	return ""
}
