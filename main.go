package main

import (
	"fmt"
	"os"

	"christhefrog/sampman/components/github"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"

	"github.com/urfave/cli/v2"
)

func main() {
	config, err := sampman.LoadConfig("sampman.json")
	if err != nil {
		util.Fatal(fmt.Sprintf("Couldn't load sampman.json (%s)", err))
	}

	release, err := github.FetchLatestRelease("pawn-lang", "compiler")

	if err != nil {
		fmt.Printf("Couldn't fetch the lastest release (%s)", err)
	} else if !util.Has(config.Compilers, release.Name) {
		fmt.Printf("\nLatest version: %s (%s)\n \n",
			release.Name, release.Published.Format("02.01.2006"))
	}

	app := &cli.App{
		Name:  "sampman",
		Usage: "A samp server manager",
		Action: func(*cli.Context) error {
			fmt.Println("Hello world!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
