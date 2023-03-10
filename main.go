package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	PROG = `contextdir`
)

var (
	// version, commit, date, builtBy are provided by goreleaser during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version %s, commit %s, built at %s by %s\n", version, commit, date, builtBy)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02T15:04:05Z07:00",
		FullTimestamp:          true})

}

func main() {
	app := &cli.App{
		Name:                 PROG,
		Version:              version,
		EnableBashCompletion: true,
		Usage:                "Utility for working with local Docker build context directories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "verbosity",
				// fixme: in a future release, urfave/cli has automatic textwrap support
				Usage: "Sets the verbosity level of the log messages printed by the program, should be\n" +
					`one of: "debug", "error", "fatal", "info", "panic", "trace", or "warn"`,
				Action: func(c *cli.Context, verbosity string) error {
					level, err := log.ParseLevel(verbosity)
					if err != nil {
						return err
					}
					log.SetLevel(level)
					return err
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "list the contents of the given context dir, honoring .dockerignore (if found)",
				ArgsUsage: "[dir]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "ignored",
						Usage: "Instead of the normal output, only list files in the given contextdir(s) which\n" +
							"match patterns in the .dockerignore file",
					},
					&cli.BoolFlag{
						Name: "detailed",
						Usage: "Format the output in JSON format with detailed information about the files (may\n" +
							"be combined with the --ignored flag)",
					},
				},
				Action: handleListCmd,
			},
			{
				Name:    "checksum",
				Aliases: []string{"sum", "hash"},
				Usage: "list the contents of the given context path(s), excluding any paths mentioned in\n" +
					".dockerignore",
				ArgsUsage: "[dir]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Format the output in JSON format",
					},
					&cli.BoolFlag{
						Name:  "detailed",
						Usage: "Report the entire scan in JSON format",
					},
				},
				Action: handleChecksumCmd,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
