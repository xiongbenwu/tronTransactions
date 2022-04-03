package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "transactions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Usage:    "GitHub API token with repo access",
				Required: true,
			},
		},
		Action: Run,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(ctx *cli.Context) error {
	return nil
}
