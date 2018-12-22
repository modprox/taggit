package updates

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/shoenig/regexplus"
)

var upRe = regexp.MustCompile(`^(?P<pkg>[\S]+)[\s](?P<old>[\S]+)[\s]\[(?P<new>[\S]+)]$`)

type Update struct {
	Package string
	Old     string
	New     string
}

func (u Update) String() string {
	return fmt.Sprintf(
		"%s %s => %s",
		u.Package,
		u.Old,
		u.New,
	)
}

func parseCmd(s string) ([]Update, error) {
	var updates []Update

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		groups := regexplus.FindNamedSubmatches(upRe, line)
		if len(groups) == 0 {
			continue
		}

		updates = append(updates, Update{
			Package: groups["pkg"],
			Old:     groups["old"],
			New:     groups["new"],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return updates, nil
}
