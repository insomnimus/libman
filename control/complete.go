package control

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"

	"github.com/insomnimus/libman/handler/cmd"
)

// completeCommand is like handler.Set.Complete but also includes user defined aliases.
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

	// complete the argument , if the command is already typed
	if strings.Contains(buf, " ") {
		command, arg := splitCmd(buf)
		h := handlers.Match(command)
		if h == nil {
			return nil
		}
		return h.Complete(command, arg)
	}

	// check handlers
	for _, h := range handlers {
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

func completeNothing(string, string) []string {
	return nil
}

func suggestPlaylist(command, arg string) []string {
	updateCache()
	pls := make([]string, 0, len(cache))
	h := handlers.Match(command)

	if arg == "" {
		// return all playlist names
		for _, p := range cache {
			// do not suggest followed playlists for the edit, save or rm commands
			if p.isFollowed && (h.Cmd == cmd.EditPlaylist ||
				h.Cmd == cmd.EditPlaylistDetails ||
				h.Cmd == cmd.SavePlaying ||
				h.Cmd == cmd.RemovePlaying) {
				continue
			}
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
		if len(pls) == 0 {
			return nil
		}
		return pls
	}

	for _, p := range cache {
		// do not suggest followed playlists for the edit, save or rm commands
		if p.isFollowed && (h.Cmd == cmd.EditPlaylist ||
			h.Cmd == cmd.EditPlaylistDetails ||
			h.Cmd == cmd.SavePlaying ||
			h.Cmd == cmd.RemovePlaying) {
			continue
		}
		if hasPrefixFold(p.Name, arg) {
			pls = append(pls,
				fmt.Sprintf("%s %s", command, p.Name))
		}
	}
	if len(pls) == 0 {
		return nil
	}
	return pls
}

// generates a completer for the second word of a command
// will return nil if the command is not a match of given command names (commands parameter)
func newWordCompleter(cand ...string) func(string, string) []string {
	return func(command, arg string) []string {
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

func dynamicCompleteFunc(vec *[]string, commands ...string) func(string, string) []string {
	return func(command, arg string) []string {
		items := make([]string, 0, len(*vec))
		for _, s := range *vec {
			if hasPrefixFold(s, arg) {
				items = append(items,
					fmt.Sprintf("%s %s", command, s))
			}
		}

		if len(items) == 0 {
			return nil
		}
		return items
	}
}

func suggestPath(buf string) []string {
	var err error
	var items []string
	if DataHome == "" || filepath.IsAbs(buf) {
		items, err = zglob.Glob(buf + "*")
	} else {
		items, err = zglob.Glob(DataHome + buf + "*")
	}
	if err != nil {
		return nil
	}
	return items
}
