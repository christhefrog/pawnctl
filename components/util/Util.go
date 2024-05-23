package util

import (
	"fmt"
	"os"
)

func Fatal(message ...any) {
	fmt.Println(message...)
	os.Exit(1)
}

func Has(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
