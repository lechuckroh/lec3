package main

import (
	"errors"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	SetLogPattern(TimeOnly)

	app := &cli.App{
		Name:  "lec",
		Usage: "",
		Commands: []cli.Command{
			{
				Name:  "ip",
				Usage: "run image processing",
				Action: func(c *cli.Context) error {
					cfg := ConfigIP{}

					cfgFilename := c.String("cfg")
					if len(cfgFilename) > 0 {
						cfg.LoadYaml(cfgFilename)
					} else {
						cfg.src.dir = c.String("src")
						cfg.dest.dir = c.String("dest")
						cfg.watch = c.Bool("watch")
					}

					cfg.Print()

					ip := ImageProcess{}
					ip.run(&cfg)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "cfg",
						Usage: "configuration file",
					},
					&cli.StringFlag{
						Name:  "src",
						Value: "./",
						Usage: "source directory",
					},
					&cli.StringFlag{
						Name:  "dest",
						Value: "./output",
						Usage: "destination directory",
					},
					&cli.BoolFlag{
						Name:  "watch",
						Usage: "watch mode",
					},
				},
			},
			{
				Name:  "conv",
				Usage: "convert images to other format (pdf, zip, ...)",
				Action: func(c *cli.Context) error {
					cfgFilename := c.String("cfg")
					src := c.String("src")
					dest := c.String("dest")

					cfg := ConfigConv{}
					if cfgFilename == "" {
						return errors.New("cfg flag value is empty")
					}
					if src != "" {
						cfg.src.filename = src
					}
					if dest != "" {
						cfg.dest.dir = dest
					}

					cfg.Print()

					conv := Convert{}
					conv.run(&cfg)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "cfg",
						Usage: "configuration file",
					},
					&cli.StringFlag{
						Name:  "src",
						Usage: "source directory",
					},
					&cli.StringFlag{
						Name:  "dest",
						Value: "./output",
						Usage: "destination filename or directory",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
