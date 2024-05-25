package main

import (
	"fmt"
	"os"

	"christhefrog/sampman/components/compiler"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"

	"github.com/urfave/cli/v2"
)

func Update(ctx *cli.Context) error {
	config, err := sampman.LoadConfig("sampman.json")
	if err != nil {
		util.Fatalf("Couldn't load sampman.json (%s)", err)
	}

	fmt.Printf("Looking for compiler updates...\n")

	release, err := compiler.FetchLatestRelease()
	if err != nil {
		util.Fatalf("Couldn't fetch latest compiler (%s)\n", err)
	}

	if config.Compilers[release.Name] == "" {
		fmt.Printf("\nA new compiler version is available: %s (%s)\nDownloading...\n\n",
			release.Name, release.Published.Format("02.01.2006"))

		err := compiler.Download(release, &config)
		if err != nil {
			util.Fatalf("Couldn't download compiler version %s (%s)", release.Name, err)
		}
	}

	fmt.Printf("Everything is up-to-date")

	return nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Check for compiler updates",
				Action:  Update,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
