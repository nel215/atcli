package setup

import (
	"errors"
	"github.com/nel215/atcli/store"
)

type Setup struct {
	configStore interface {
		Save(*store.Config) error
	}
}

func New() *Setup {
	return &Setup{store.NewConfigStore()}
}

func (s *Setup) Execute(contestUrl string) error {
	if contestUrl == "" {
		return errors.New("contestUrl is required")
	}
	config := &store.Config{contestUrl}
	err := s.configStore.Save(config)
	if err != nil {
		return err
	}
	return nil
}
