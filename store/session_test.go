package store

import (
	"github.com/nel215/atcmd/login"
	"io/ioutil"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	sess := &login.Session{
		"_session",
		"_issue_time",
		"_kick_id",
		"_user_id",
	}
	ss := NewSessionStore()
	err := ss.Save(sess)
	if err != nil {
		t.Fatalf("save failed")
	}

	f, err := os.Open("./.session.json")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("%s", err)
	}
	expected := `{"Session":"_session","IssueTime":"_issue_time","KickId":"_kick_id","UserID":"_user_id"}`
	if string(bs) != expected {
		t.Fatal("expected %s. bot %s", expected, string(bs))
	}
}
