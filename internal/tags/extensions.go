package tags

import (
	"flag"
	"strings"
)

// Extensions represent the two kinds of extra "labels" that can be attached to
// a semantic version tag, described in https://semver.org/.
type Extensions struct {
	PreRelease   string
	BuildMetadta string
}

// ExtractExtensions will extract the pre-release and/or build-metadata information from
// a FlagSet, if those pieces of information were provided.
//
// The pre-release information is the 0th positional argument in the context of
// a sub-command.
//
// The build-metadata information is passed in via the "-meta" flag.
func ExtractExtensions(fs *flag.FlagSet) Extensions {
	var (
		preRelease    string
		buildMetadata string
	)

	args := fs.Args()
	if len(args) > 0 {
		preRelease = args[0]
	}

	buildMetadata = fs.Lookup("meta").Value.String()

	return Extensions{
		PreRelease:   clean(preRelease),
		BuildMetadta: clean(buildMetadata),
	}
}

func clean(orig string) string {
	noDash := strings.TrimPrefix(orig, "-")
	noPlus := strings.TrimPrefix(noDash, "+")
	return noPlus
}
