package main

import (
	"fmt"
	"os"

	"christhefrog/sampman/components/compiler"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"

	"github.com/urfave/cli/v2"
)

func main() {
	config, err := sampman.LoadConfig("sampman.json")
	if err != nil {
		util.Fatalf("Couldn't load sampman.json (%s)", err)
	}

	release, err := compiler.FetchLatestCompiler()
	if err != nil {
		fmt.Printf("Couldn't fetch latest compiler (%s)", err)
	}

	if !config.GetCompiler(release.Name).IsInstalled() {
		fmt.Printf("\nA new compiler version is available: %s (%s)\nDownloading...\n\n",
			release.Name, release.Published.Format("02.01.2006"))

		err := compiler.Download(release, &config)
		if err != nil {
			fmt.Printf("Couldn't download compiler version %s (%s)", release.Name, err)
		}
	}

	app := &cli.App{
		Name:  "sampman",
		Usage: "A samp server manager",
		Action: func(*cli.Context) error {
			fmt.Println("Everything up-to-date!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
