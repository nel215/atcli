package api

import (
	"fmt"
	"github.com/nel215/atcli/session"
	"github.com/nel215/atcli/store"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type History struct {
	sessionStore interface {
		Load() (*session.Session, error)
	}
	configStore interface {
		Load() (*store.Config, error)
	}
}

func NewHistory() *History {
	return &History{
		sessionStore: store.NewSessionStore(),
		configStore:  store.NewConfigStore(),
	}
}

func (h *History) Execute() error {
	config, err := h.configStore.Load()
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

	sess, err := h.sessionStore.Load()
	if err != nil {
		return err
	}
	cookies := sess.GetSessionCookies()
	jar.SetCookies(u, cookies)
	http.DefaultClient.Jar = jar

	resp, err := http.Get(contestUrl + "/submissions/me")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	texts := []string{}
	fmt.Print("# History\n\n")
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data != "td" {
				continue
			}
			for {
				z.Next()
				t = z.Token()
				if t.Data == "td" {
					break
				}
				if t.Type != html.TextToken {
					continue
				}
				texts = append(texts, t.Data)
			}

		case tt == html.EndTagToken:
			t := z.Token()

			if t.Data == "tr" {
				if len(texts) > 0 {
					fmt.Printf("- %s\n", strings.Join(texts, ", "))
					texts = []string{}
				}
			}
		}
	}

	return nil
}
