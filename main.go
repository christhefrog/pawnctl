package main

import (
	"fmt"
	"os"

	"christhefrog/pawnctl/components/compiler"
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/project"
	"christhefrog/pawnctl/components/util"

	"github.com/urfave/cli/v2"
)

func Update(ctx *cli.Context) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load pawnctl.json (%s)", err)
	}

	fmt.Printf("Looking for compiler updates...\n")

	release, err := compiler.FetchLatestRelease()
	if err != nil {
		util.Fatalf("Couldn't fetch latest compiler (%s)\n", err)
	}

	if !config.IsCompilerInstalled(release.Name) {
		fmt.Printf("\nA new compiler version is available: %s (%s)\nDownloading...\n\n",
			release.Name, release.Published.Format("02.01.2006"))

		err := compiler.Download(release)
		if err != nil {
			util.Fatalf("Couldn't download compiler version %s (%s)", release.Name, err)
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

func Init(ctx *cli.Context) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load pawnctl.json (%s)", err)
	}

	proj, err := project.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load project pawnctl.json (%s)", err)
	}

	if proj.CompilerVersion != "" {
		util.Fatalf("Project is already initialized")
	}

	version := ""
	fmt.Printf("\nCompiler version (leave blank for latest)\n%v\n> ", config.ListCompilers())
	fmt.Scanln(&version)

	if version == "" {
		version = "latest"
	}

	source := ""
	fmt.Print("\nSource file (leave blank for gamemodes\\gamemode.pwn)\n> ")
	fmt.Scanln(&source)

	if source == "" {
		source = "gamemodes\\gamemode.pwn"
	}

	include := ""
	fmt.Print("\nInclude path (leave blank for qawno\\include)\n> ")
	fmt.Scanln(&include)

	if include == "" {
		include = "qawno\\include"
	}

	proj.CompilerVersion = version
	proj.Sources = []string{source}
	proj.Includes = []string{include}

	proj.Save()

	return nil
}

func main() {
	app := &cli.App{
		Name:  "pawnctl",
		Usage: "A pawn installation manager",
		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update the compiler",
				Action:  Update,
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize a project",
				Action:  Init,
			},
			{
				Name:    "compile",
				Aliases: []string{"c"},
				Usage:   "Compile a project",
				Action:  Compile,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
