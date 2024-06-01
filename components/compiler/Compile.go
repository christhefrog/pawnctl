package compiler

import (
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/util"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Compile(source string, args ...string) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load pawnctl.json (%s)", err)
	}

	_, compiler := config.GetLatestCompiler()

	if compiler == "" {
		return errors.New("couldn't find the latest compiler, use `sampman u`")
	}

	arg := make([]string, 0)
	arg = append(arg, source)
	arg = append(arg, args...)

	cmd := exec.Command(compiler, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%s %s\n", compiler, strings.Join(arg[:], " "))

	start := time.Now()

	_ = cmd.Run()

	fmt.Printf("...took %s", time.Since(start))

	return nil
}
