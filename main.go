package main

import (
	"fmt"
	"os"
	"time"

	"christhefrog/pawnctl/components/compiler"
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/project"
	"christhefrog/pawnctl/components/util"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func Update(ctx *cli.Context) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	color.Gray.Printf("Looking for compiler updates...\n")

	release, err := compiler.FetchLatestRelease()
	if err != nil {
		util.Fatalf("Couldn't fetch latest compiler (%s)\n", err)
	}

	if !config.IsCompilerInstalled(release.Name) {
		color.Blue.Printf("\nA new compiler version is available: %s (%s)",
			release.Name, release.Published.Format("02.01.2006"))
		fmt.Print("\nDownloading...\n\n")

		err := compiler.Download(release)
		if err != nil {
			util.Fatalf("Couldn't download compiler version %s (%s)", release.Name, err)
		}
	}

	color.Green.Printf("Everything is up-to-date")

	return nil
}

func Compile(ctx *cli.Context) error {
	profile := ctx.Args().First()

	compiler.Compile(profile)

	return nil
}

func Watch(ctx *cli.Context) error {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	proj, err := project.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load project config (%s)", err)
	}

	if len(proj.Profiles) < 1 {
		util.Fatalf("Project config not found, use `pawnctl i`")
	}

	profile := ctx.Args().First()

	prof, ok := proj.Profiles[profile]
	if !ok {
		util.Fatalf("Profile %s doesn't exist in current project", profile)
	}

	fmt.Print("\033[H\033[2J") // Clear screen
	compiler.Compile(profile)
	color.Gray.Print("Watching for changes...\n")

	watcher.Add(prof.Input)
	for _, v := range prof.Includes {
		watcher.Add(v)
	}

	go func() {
		var (
			timer     *time.Timer
			lastEvent fsnotify.Event
		)
		timer = time.NewTimer(time.Millisecond)
		<-timer.C // timer should be expired at first
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				lastEvent = event
				timer.Reset(time.Millisecond * 200)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				util.Fatalf("Error watching for changes (%s)", err)
			case <-timer.C:
				if lastEvent.Op&fsnotify.Write == fsnotify.Write {
					fmt.Print("\033[H\033[2J") // Clear screen
					compiler.Compile(profile)
					color.Gray.Print("Watching for changes...\n")
				}
				if err != nil {
					util.Fatalf("Error watching for changes (%s)", err)
				}
			}

		}
	}()

	<-make(chan struct{})
	return nil
}

func Init(ctx *cli.Context) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	if len(config.ListCompilers()) < 1 {
		util.Fatalf("No compilers found, use `pawnctl u`")
	}

	proj, err := project.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load project config (%s)", err)
	}

	if len(proj.Profiles) > 0 {
		util.Fatalf("Project is already initialized")
	}

	version := ""
	fmt.Printf("\nCompiler version ")
	color.Gray.Printf("(leave blank for latest)\n%v\n", config.ListCompilers())
	fmt.Print("> ")
	fmt.Scanln(&version)

	if version == "" {
		version = config.Compilers["latest"]
	}

	source := ""
	fmt.Print("\nSource file ")
	color.Gray.Print("(leave blank for gamemodes\\gamemode.pwn)\n")
	fmt.Print("> ")
	fmt.Scanln(&source)

	if source == "" {
		source = "gamemodes\\gamemode.pwn"
	}

	output := ""
	fmt.Print("\nOutput ")
	color.Gray.Print("(leave blank for gamemodes\\gamemode.amx)\n")
	fmt.Print("> ")
	fmt.Scanln(&output)

	if output == "" {
		output = "gamemodes\\gamemode.amx"
	}

	include := ""
	fmt.Print("\nInclude path ")
	color.Gray.Print("(leave blank for qawno\\include)\n")
	fmt.Print("> ")
	fmt.Scanln(&include)

	if include == "" {
		include = "qawno\\include"
	}

	proj.Profiles[""] = project.Profile{
		CompilerVersion: version,
		Input:           source,
		Output:          output,
		Includes:        []string{include},
		Args:            []string{"-d3", "-Z-", "-;+", "-(+", "-\\", "-t4"},
	}

	proj.Profiles["release"] = project.Profile{
		CompilerVersion: version,
		Input:           source,
		Output:          output,
		Includes:        []string{include},
		Args:            []string{"-d0", "-O2", "-Z-", "-;+", "-(+", "-\\", "-t4"},
	}

	proj.Save()

	color.Green.Print("\nYou can now use:\n")
	color.Gray.Print("pawnctl c ")
	color.Blue.Print("\t\tto build a debug version\n")
	color.Gray.Print("pawnctl c release ")
	color.Blue.Print("\tto build a release version\n")
	color.Gray.Print("pawnctl w (release) ")
	color.Blue.Print("\tto build a debug/release version every time a file changes\n \n")

	color.Gray.Print("If you want to create a new profile check out pawnctl.json")

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
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch for changes in a file and compile",
				Action:  Watch,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
