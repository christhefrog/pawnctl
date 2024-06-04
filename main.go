package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"christhefrog/pawnctl/components/compiler"
	"christhefrog/pawnctl/components/github"
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

	var release github.Release

	name := ctx.Args().First()
	if name == "" || name == "latest" {
		color.Gray.Printf("Looking for lastest compiler...\n")

		release, err = compiler.FetchLatestRelease()
		if err != nil {
			util.Fatalf("Couldn't fetch latest compiler (%s)\n", err)
		}
	} else {
		color.Gray.Printf("Looking for compiler version %s...\n", name)

		release, err = compiler.FetchRelease(name)
		if err != nil {
			util.Fatalf("Couldn't fetch compiler version %s (%s)\n", name, err)
		}
	}

	if !config.IsCompilerInstalled(release.Name) {
		color.Blue.Printf("\nCompiler version %s is available (%s)",
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

	parent := filepath.Dir(prof.Input)

	fmt.Print("\033[H\033[2J") // Clear screen

	// See christhefrog/pawnctl#1
	color.Red.Printf("\nPlease note that for the time, watch only looks for changes in %s\\**.\n", parent)
	color.Gray.Print("It means that other directories specified in profile.includes (ex. qawno\\include) aren't scanned.\n \n")

	compiler.Compile(profile)
	color.Gray.Print("Watching for changes...\n")

	filepath.Walk(parent, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		watcher.Add(path)
		return nil
	})

	var (
		timer     *time.Timer
		lastEvent fsnotify.Event
	)
	timer = time.NewTimer(time.Millisecond)
	<-timer.C // timer should be expired at first
	for {
		select {
		case event, ok := <-watcher.Events:
			ext := filepath.Ext(event.Name)
			if !ok || (ext != ".pwn" && ext != ".inc" && ext != ".p" && ext != ".pawn" && ext != "") {
				continue
			}
			lastEvent = event
			// Sometimes watch fires off twice, timer helps with that
			timer.Reset(time.Millisecond * 200)
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			util.Fatalf("Error watching for changes (%s)", err)
		case <-timer.C:
			ext := filepath.Ext(lastEvent.Name)
			if lastEvent.Op.Has(fsnotify.Write) && ext != "" {
				fmt.Print("\033[H\033[2J") // Clear screen
				color.Green.Println(lastEvent.Name, "changed")
				compiler.Compile(profile)
				color.Gray.Print("Watching for changes...\n")
			}
			// Is directory?
			if lastEvent.Op.Has(fsnotify.Create) && ext == "" {
				watcher.Add(lastEvent.Name)
			}
			if err != nil {
				util.Fatalf("Error watching for changes (%s)", err)
			}
		}

	}
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

	if version == "" || version == "latest" {
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
		Args:            []string{"-d0", "-O1", "-Z-", "-;+", "-(+", "-\\", "-t4"},
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
				// Action:  Watch,
				Action: Watch,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		util.Fatal(err)
	}
}
