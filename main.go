package main

import (
	github "christhefrog/pawnman/components"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pawnman",
		Usage: "A pawn installation and configuration manager",
		Action: func(*cli.Context) error {
			release := github.FetchLatestRelease("pawn-lang", "compiler")

			fmt.Printf("Latest release: %s (%s)\n", release.Name, release.Published.Format("02.01.2006"))
			for _, v := range release.Assets {
				fmt.Println(v.Name)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
