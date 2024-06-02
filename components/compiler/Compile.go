package compiler

import (
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/project"
	"christhefrog/pawnctl/components/util"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gookit/color"
)

func CompileFileWithDefaults(file string) {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	compiler := config.Compilers[config.Compilers["latest"]]
	if compiler == "" {
		util.Fatalf("No compilers found, use `pawnctl u`")
	}

	cmd := exec.Command(compiler, file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%s %s\n", compiler, file)

	start := time.Now()

	_ = cmd.Run()

	fmt.Printf("...took %s", time.Since(start))
}

func Compile(profile string) {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	proj, err := project.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load project config (%s)", err)
	}

	if len(proj.Profiles) < 1 {
		util.Fatalf("Project config not found, use `pawnctl i`")
	}

	prof, ok := proj.Profiles[profile]

	if !ok {
		util.Fatalf("Profile %s doesn't exist in current project", profile)
	}

	compiler := config.Compilers[prof.CompilerVersion]
	if prof.CompilerVersion == "latest" {
		compiler = config.Compilers[compiler]
	}

	if compiler == "" {
		util.Fatalf("Couldn't find the compiler version %s, use `pawnctl u`", prof.CompilerVersion, prof.CompilerVersion)
	}

	args := make([]string, 0)
	args = append(args, prof.Input)
	args = append(args, fmt.Sprint("-o", prof.Output))

	for _, v := range prof.Includes {
		args = append(args, fmt.Sprint("-i", v))
	}

	args = append(args, prof.Args...)

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	color.Gray.Printf("%s %s\n", compiler, strings.Join(args[:], " "))

	start := time.Now()

	_ = cmd.Run()

	// file := filepath.Base(proj.Input)
	// os.Rename(fmt.Sprint(strings.TrimSuffix(file, filepath.Ext(file)), ".amx"), proj.Output)

	color.Gray.Printf("...took %s\n", time.Since(start))
}
