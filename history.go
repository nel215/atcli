package main

import (
	"github.com/nel215/atcli/session"
	"io/ioutil"
	"log"
	"net/http"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	return nil
}
