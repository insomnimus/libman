package control

import (
	"fmt"
	"libman/handler/cmd"
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
		// check if the line is a playlist command
		if (h.Cmd == cmd.PlayUserPlaylist || h.Cmd == cmd.DeletePlaylist || h.Cmd == cmd.EditPlaylist) &&
			h.HasPrefix(buf) {
			return suggestPlaylist(buf)
		}

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

func suggestPlaylist(buf string) []string {
	updateCache()
	pls := make([]string, 0, len(*cache))
	buf = strings.ToLower(buf)
	command, name := splitCmd(buf)

	if name == "" {
		// return all playlist names
		for _, p := range *cache {
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
		return pls
	}

	for _, p := range *cache {
		if strings.HasPrefix(strings.ToLower(p.Name), name) {
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
	}
	return pls
}
