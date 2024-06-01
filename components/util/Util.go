package util

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

func Fatal(message ...any) {
	fmt.Println(message...)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	color.Red.Printf(fmt.Sprint(format, "\n"), args...)
	os.Exit(1)
}
