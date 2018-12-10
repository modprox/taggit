package cli

import (
	"fmt"
	"os"
)

func Usage(code int) {
	fmt.Println("usage: taggit [list, zero, patch, minor, major]")
	os.Exit(code)
}
