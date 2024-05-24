package util

import (
	"fmt"
	"os"
)

func Fatal(message ...any) {
	fmt.Println(message...)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	fmt.Printf(fmt.Sprint(format, "\n"), args...)
	os.Exit(1)
}
