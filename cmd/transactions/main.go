package main

import (
	"log"
	"os"
	"transactions/pkg/client"
	"transactions/pkg/server"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "transactions",
		Usage: "TRON transactions server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Usage:    "Main account address",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "API server URL",
				Value: "https://api.trongrid.io/v1/",
			},
		},
		Action: Run,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(ctx *cli.Context) error {
	tr, err := client.NewClient(ctx.String("address"), ctx.String("url"))
	if err != nil {
		return err
	}
	err = tr.Run()
	if err != nil {
		return err
	}
	server := server.Server{Client: tr}
	server.ListenAndServe()
	return nil
}
