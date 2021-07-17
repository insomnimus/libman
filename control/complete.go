package control

import (
	"strings"
)

func completeCommand(buf string) (c []string) {
	buf = strings.ToLower(buf)
	// first check aliases
	for _, a := range userAliases.Inner() {
		if strings.HasPrefix(strings.ToLower(a.Left), buf) {
			c = append(c, a.Left)
		}
	}

	// check handlers
	for _, h := range handlers {
		if strings.HasPrefix(h.Name, buf) {
			c = append(c, h.Name)
			continue
		}

		for _, a := range h.Aliases {
			if strings.HasPrefix(a, buf) {
				c = append(c, a)
				break
			}
		}
	}

	return
}

func completeBool(buf string) (c []string) {
	buf = strings.ToLower(buf)
	if strings.HasPrefix("yes", buf) {
		c = []string{"yes"}
	} else if strings.HasPrefix("no", buf) {
		c = []string{"no"}
	}
	return c

}

func completeNothing(string) []string {
	return nil
}
