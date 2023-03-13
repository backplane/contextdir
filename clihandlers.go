package main

import (
	"fmt"

	"github.com/backplane/contextdir/scan"
	"github.com/backplane/contextdir/utils"

	"github.com/urfave/cli/v2"
)

// checksumReport represents the fields returned by the "checksum" command in "--json" mode
type checksumReport struct {
	Path     string
	Checksum string
}

func handleListCmd(c *cli.Context) error {
	onlyIgnored := c.Bool("ignored")
	detailed := c.Bool("detailed")

	var path string
	switch c.NArg() {
	case 0:
		path = "."
	case 1:
		path = c.Args().First()
	default:
		return fmt.Errorf("only one directory argument is allowed")
	}

	scan, err := scan.ScanContextDir(path)
	if err != nil {
		return err
	}

	if onlyIgnored {
		if detailed {
			if err := utils.WriteJSON(scan.IgnoredFiles); err != nil {
				return err
			}
			return nil
		}
		// if not hyper-verbose json, just produce a plain file list
		for _, f := range scan.IgnoredFiles {
			fmt.Println(f.Stats.Path)
		}
		return nil
	}

	if detailed {
		if err := utils.WriteJSON(scan.Dir.Files); err != nil {
			return err
		}
		return nil
	}
	for _, f := range scan.Dir.Files {
		fmt.Println(f.Stats.Path)
	}
	return nil
}

func handleChecksumCmd(c *cli.Context) error {
	jsonOutput := c.Bool("json")
	detailed := c.Bool("detailed")

	var path string
	switch c.NArg() {
	case 0:
		path = "."
	case 1:
		path = c.Args().First()
	default:
		return fmt.Errorf("only one directory argument is allowed")
	}

	scan, err := scan.ScanContextDir(path)
	if err != nil {
		return err
	}

	if detailed {
		// write the whole scan
		if err := utils.WriteJSON(scan.Dir); err != nil {
			return err
		}
		return nil
	}

	report := &checksumReport{
		Path:     scan.Dir.Path,
		Checksum: scan.Dir.GetHashString(),
	}
	if jsonOutput {
		if err := utils.WriteJSON(report); err != nil {
			return err
		}
		return nil
	}

	fmt.Println(report.Checksum)
	return nil
}
