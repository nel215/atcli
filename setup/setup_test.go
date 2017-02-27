package setup

import (
	"github.com/nel215/atcli/store"
	"testing"
)

type TestingConfigStore struct {
	config *store.Config
}

func (tcs *TestingConfigStore) Save(config *store.Config) error {
	tcs.config = config
	return nil
}

func TestExecute(t *testing.T) {
	u := "https://localhost"
	s := &Setup{&TestingConfigStore{}}
	err := s.Execute(u)
	if err != nil {
		t.Fatalf("setup failed\n")
	}
	if actual := s.configStore.(*TestingConfigStore).config.ContestUrl; actual != u {
		t.Fatalf("expected url is %s. got %s\n", u, actual)
	}
}
