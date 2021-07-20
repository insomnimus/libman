package control

import (
	"fmt"
	"strings"

	"github.com/insomnimus/libman/handler/cmd"
)

func completeCommand(buf string) (c []string) {
	buf = strings.TrimPrefix(buf, " ")
	// first check aliases
	for _, a := range userAliases.Inner() {
		if hasPrefixFold(a.Left, buf) {
			c = append(c, a.Left)
		}
	}
	// do not include handler aliases if buf has no text
	if buf == "" {
		return append(c, handlers.Names()...)
	}
	hasSpace := strings.Contains(buf, " ")
	// check handlers
	for _, h := range handlers {
		if hasSpace {
			candidates := h.Complete(buf)
			if candidates != nil {
				return candidates
			}
			continue
		}
		// complete the command itself
		if hasPrefixFold(h.Name, buf) {
			c = append(c, h.Name)
		}
		for _, a := range h.Aliases {
			if hasPrefixFold(a, buf) {
				c = append(c, a)
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
	command, name := splitCmd(buf)
	if command == "" {
		return nil
	}
	// do nothing if the command is not a playlist command
	h := handlers.Match(command)
	if h == nil ||
		!(h.Cmd == cmd.PlayUserPlaylist ||
			h.Cmd == cmd.EditPlaylistDetails ||
			h.Cmd == cmd.SavePlaying ||
			h.Cmd == cmd.RemovePlaying ||
			h.Cmd == cmd.EditPlaylist ||
			h.Cmd == cmd.DeletePlaylist) {
		return nil
	}

	if name == "" {
		// return all playlist names
		for _, p := range *cache {
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
		return pls
	}

	for _, p := range *cache {
		if hasPrefixFold(p.Name, name) {
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
	}
	return pls
}

// generates a completer for the second word of a command
// will return nil if the command is not a match of given command names (commands parameter)
func newWordCompleter(cand []string, commands ...string) func(string) []string {
	return func(buf string) []string {
		buf = strings.TrimPrefix(buf, " ") // ignore leading space
		command, arg := splitCmd(buf)
		isMatch := false
		for _, c := range commands {
			if strings.EqualFold(command, c) {
				isMatch = true
				break
			}
		}
		if !isMatch {
			return nil
		}
		items := make([]string, 0, len(cand))
		for _, s := range cand {
			if hasPrefixFold(s, arg) {
				items = append(items,
					fmt.Sprintf("%s %s", command, s))
			}
		}
		// we want to return nil in case of no match
		if len(items) == 0 {
			return nil
		}
		return items
	}
}
