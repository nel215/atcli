package api

import (
	"fmt"
	"github.com/nel215/atcli/session"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func DescribeHistory(sess *session.Session) error {
	req, err := http.NewRequest(http.MethodGet, "https://practice.contest.atcoder.jp/submissions/me", nil)
	if err != nil {
		return err
	}
	sess.AddSessionCookies(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	texts := make([]string, 0)
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
				}
			}
		}
	}

	return nil
}