package compiler

import (
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/project"
	"christhefrog/pawnctl/components/util"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func Compile() {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	proj, err := project.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load project config (%s)", err)
	}

	if proj.CompilerVersion == "" {
		util.Fatalf("Project config not found, use `pawnctl i`")
	}

	compiler := config.Compilers[proj.CompilerVersion]
	if proj.CompilerVersion == "latest" {
		compiler = config.Compilers[compiler]
	}

	if compiler == "" {
		util.Fatalf("Couldn't find the compiler version %s, use `pawnctl u`", proj.CompilerVersion, proj.CompilerVersion)
	}

	args := make([]string, 0)
	args = append(args, proj.Input)
	//args = append(args, fmt.Sprint("-D", proj.OutputDir))

	for _, v := range proj.Includes {
		args = append(args, fmt.Sprint("-i", v))
	}

	args = append(args, proj.Args...)

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	color.Gray.Printf("%s %s\n", compiler, strings.Join(args[:], " "))

	start := time.Now()

	_ = cmd.Run()

	file := filepath.Base(proj.Input)
	os.Rename(fmt.Sprint(strings.TrimSuffix(file, filepath.Ext(file)), ".amx"), proj.Output)

	color.Gray.Printf("...took %s\n", time.Since(start))
}
