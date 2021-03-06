package main

import (
	"errors"
	"github.com/nel215/atcli/api"
	"github.com/nel215/atcli/login"
	"github.com/nel215/atcli/problem"
	"github.com/nel215/atcli/setup"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "atcli"
	app.Commands = []cli.Command{
		{
			Name: "setup",
			Action: func(c *cli.Context) error {
				contestUrl := c.String("u")
				s := setup.New()
				err := s.Execute(contestUrl)
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "u"},
			},
		},
		{
			Name: "login",
			Action: func(c *cli.Context) error {
				user := c.String("u")
				password := c.String("p")
				if user == "" {
					return errors.New("user(-u) is required")
				}
				if password == "" {
					return errors.New("password(-p) is required")
				}

				log.Printf("try logging in by %s...\n", user)
				l := login.New()
				err := l.Submit(user, password)
				if err != nil {
					return err
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "u"},
				cli.StringFlag{Name: "p"},
			},
		},
		{
			Name:  "history",
			Usage: "describe submission history",
			Action: func(c *cli.Context) error {
				err := api.NewHistory().Execute()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "problem",
			Usage: "describe problem list",
			Action: func(c *cli.Context) error {
				p := problem.New()
				err := p.Execute()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "language",
			Usage: "describe language list of specified problem",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "submit",
			Usage: "submit source code",
			Action: func(c *cli.Context) error {
				problemId := c.Int64("p")
				languageId := c.Int64("l")
				sourceCodePath := c.String("s")

				s, err := api.NewSubmit(problemId, languageId, sourceCodePath)
				if err != nil {
					return err
				}
				err = s.Execute()
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.Int64Flag{Name: "p"},
				cli.Int64Flag{Name: "l"},
				cli.StringFlag{Name: "s"},
			},
		},
	}
	app.Run(os.Args)
}
