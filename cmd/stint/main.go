package main

import (
	"os"

	"github.com/umayr/stint"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "stint"
	app.Usage = "a super tiny worker that runs in background and fetches torrent information from RSS feeds"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level, l",
			Usage: "set logging level, available options are debug, info, warn, error, and fatal",
		},
		cli.StringFlag{
			Name:  "config-file, f",
			Usage: "path to configuration file",
		},
	}

	app.Action = func(c *cli.Context) error {
		return stint.Do(c.String("config-file"), c.String("log-level"))
	}

	app.Run(os.Args)
}
