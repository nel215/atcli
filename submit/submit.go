package submit

import (
	"errors"
	"fmt"
	"github.com/nel215/atcli/login"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func createForm(problemId int64, languageId int64, sourceCode []byte, csrfToken string) url.Values {
	data := url.Values{}
	data.Set("__session", csrfToken)
	data.Set("task_id", fmt.Sprintf("%d", problemId))
	data.Set(fmt.Sprintf("language_id_%d", problemId), fmt.Sprintf("%d", languageId))
	data.Set("source_code", string(sourceCode))
	return data
}

func extractCSRFToken(sess *login.Session) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://practice.contest.atcoder.jp/submit", nil)

	resp, err := http.DefaultClient.Do(req)
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

func Submit(sess *login.Session, problemId int64, languageId int64, sourceCode []byte) error {
	jar, err := cookiejar.New(nil)
	cookies := sess.GetSessionCookies()
	u, err := url.Parse("https://practice.contest.atcoder.jp")
	if err != nil {
		return err
	}
	jar.SetCookies(u, cookies)
	http.DefaultClient.Jar = jar
	csrfToken, err := extractCSRFToken(sess)
	if err != nil {
		return err
	}
	data := createForm(problemId, languageId, sourceCode, csrfToken)
	submitUrl := fmt.Sprintf("https://practice.contest.atcoder.jp/submit?task_id=%d", problemId)
	resp, err := http.DefaultClient.PostForm(submitUrl, data)
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
