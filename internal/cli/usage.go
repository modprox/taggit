package cli

import (
	"fmt"
	"os"
	"strings"
)

const usageText = `
ABOUT

	taggit - create and manage Go module version tags

USAGE

	list
		show valid version tags in the repository

	zero
		create an initial v0.0.0 tag

	patch [pre-release]
		create an incremented patch version tag (vA.B.C+1)
		if pre-release is specified, it is appended as an extension (vA.B.C-rc1)

	minor [pre-release]
		create an incremented minor version tag (vA.B+1.C)
		if pre-release is specified, it is appended as an extension (vA.B.C-rc1)

	major [pre-release]
		create an incremented major version tag (v(A+1.B.C)
		if pre-release is specified, it is appended as an extension (vA.B.C-rc1)

	updates [apply]
		list newer versions available for outdated modules in use
		if apply is specified, automatically update the go.mod file
`

func Usage(code int) {
	fmt.Println(strings.TrimSpace(usageText))
	os.Exit(code)
}
