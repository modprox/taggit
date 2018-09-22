package cli

import (
	"fmt"
	"os"
)

func Usage(code int) {
	fmt.Fprintf(os.Stderr, "usage: taggit [list, zero, patch, minor, major]\n")
	os.Exit(code)
}
