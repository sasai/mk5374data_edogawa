package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mk5374data_edogawa"
	app.Usage = "5374.jp 江戸川区版データ作成ツール"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "target",
			Usage:  "Show content for target.csv",
			Action: RunTarget,
		},
		{
			Name:   "areadays",
			Usage:  "Show content for areadays.csv",
			Action: RunAreadays,
		},
	}
	app.Run(os.Args)
}

func RunTarget(c *cli.Context) error {
	err := mkTarget()
  return err
}

func RunAreadays(c *cli.Context) error {
	err := mkAreadays()
  return err
}
