package main

import (
	"context"
	"github.com/knowhunger/ortoo/pkg/log"
	"github.com/knowhunger/ortoo/server/server"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Usage:    "server configuration file in JSON format",
				Required: true,
			},
		},

		Action: func(c *cli.Context) error {
			confFile := c.String("conf")

			conf, err := server.LoadOrtooServerConfig(confFile)
			if err != nil {
				os.Exit(1)
			}
			log.Logger.Infof("Config: %#v", conf)
			svr, err := server.NewOrtooServer(context.Background(), conf)
			if err != nil {
				_ = log.OrtooError(err)
				os.Exit(1)
			}
			go func() {
				if err := svr.Start(); err != nil {
					_ = log.OrtooError(err)
					os.Exit(1)
				}
			}()
			os.Exit(svr.HandleSignals())
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		_ = log.OrtooError(err)
	}
}
