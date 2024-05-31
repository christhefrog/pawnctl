package main

import (
	"fmt"
	"os"

	"christhefrog/sampman/components/compiler"
	"christhefrog/sampman/components/openmp"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"

	"github.com/urfave/cli/v2"
)

func Update(ctx *cli.Context) error {
	err := sampman.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load sampman.json (%s)", err)
	}

	fmt.Printf("Looking for compiler updates...\n")

	release, err := compiler.FetchLatestRelease()
	if err != nil {
		util.Fatalf("Couldn't fetch latest compiler (%s)\n", err)
	}

	if !sampman.IsCompilerInstalled(release.Name) {
		fmt.Printf("\nA new compiler version is available: %s (%s)\nDownloading...\n\n",
			release.Name, release.Published.Format("02.01.2006"))

		err := compiler.Download(release)
		if err != nil {
			util.Fatalf("Couldn't download compiler version %s (%s)", release.Name, err)
		}
	}

	fmt.Printf("Looking for server updates...\n")

	release, err = openmp.FetchLatestServer()
	if err != nil {
		util.Fatalf("Couldn't fetch latest server (%s)\n", err)
	}

	if !sampman.IsServerInstalled(release.Name) {
		fmt.Printf("\nA new server version is available: %s (%s)\nDownloading...\n\n",
			release.Name, release.Published.Format("02.01.2006"))

		err = openmp.Download(release)
		if err != nil {
			util.Fatalf("Couldn't download server version %s (%s)", release.Name, err)
		}
	}

	fmt.Printf("Everything is up-to-date")

	return nil
}

func Compile(ctx *cli.Context) error {
	filename := ctx.Args().First()
	if filename == "" {
		filename = "gamemode.pwn"
	}

	if _, err := os.Stat(filename); err != nil {
		util.Fatalf("Couldnt find %s", filename)
	}

	err := compiler.Compile(filename)
	if err != nil {
		util.Fatalf("Couldn't compile %s (%s)", filename, err)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "sampman",
		Usage: "A samp server manager",
		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Check for compiler updates",
				Action:  Update,
			},
			{
				Name:    "compile",
				Aliases: []string{"c"},
				Usage:   "Compile a specified pawn source (gamemode.pwn by default)",
				Action:  Compile,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
