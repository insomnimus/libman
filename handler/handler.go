package handler

import (
	"fmt"
	"strings"

	"github.com/insomnimus/libman/util"
)

type Handler struct {
	Name     string
	Aliases  []string
	Cmd      uint8
	Help     string
	About    string
	Usage    string
	Run      func(string) error
	Complete func(string) []string
}

type Set []Handler

func (s *Set) Find(c uint8) *Handler {
	for i, h := range *s {
		if h.Cmd == c {
			return &(*s)[i]
		}
	}
	return nil
}

func (h Handler) String() string {
	if len(h.Aliases) == 0 {
		return fmt.Sprintf("#%s\n  %s", h.Name, h.About)
	}
	return fmt.Sprintf("#%s [aliases: %s]\n  %s", h.Name, strings.Join(h.Aliases, ", "), h.About)
}

func (h Handler) GoString() string {
	if len(h.Aliases) == 0 {
		return fmt.Sprintf("#%s\nusage:\n  %s\n  %s", h.Name, h.Usage, h.Help)
	}

	return fmt.Sprintf(`#%s
aliases:
  %s
usage:
  %s
  %s`, h.Name, strings.Join(h.Aliases, ", "), h.Usage, h.Help)
}

func (s Set) ShowUsage(c uint8) {
	if h := s.Find(c); h != nil {
		fmt.Printf("usage:\n  %s\n", h.Usage)
	}
}

func (s *Set) Match(cmd string) *Handler {
	for i, h := range *s {
		if h.Matches(cmd) {
			return &(*s)[i]
		}
	}
	return nil
}

func (h *Handler) Matches(s string) bool {
	if strings.EqualFold(h.Name, s) {
		return true
	}
	for _, a := range h.Aliases {
		if strings.EqualFold(a, s) {
			return true
		}
	}
	return false
}

func (s *Set) Len() int {
	return len(*s)
}

func (h *Handler) HasPrefix(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	if util.HasPrefixFold(h.Name, s) {
		return true
	}

	for _, a := range h.Aliases {
		if util.HasPrefixFold(s, a) {
			return true
		}
	}

	return false
}

func (s Set) CommandsAndAliases() []string {
	items := make([]string, 0, len(s))
	for _, h := range s {
		items = append(items, h.Name)
		items = append(items, h.Aliases...)
	}
	return items
}

func (s Set) Names() (names []string) {
	for _, h := range s {
		names = append(names, h.Name)
	}
	return
}

func (s Set) Complete(buf string) []string {
	buf = strings.TrimPrefix(buf, " ")
	// do not include handler aliases if buf has no text
	if buf == "" {
		return s.Names()
	}
	c := make([]string, 0, len(s))
	hasSpace := strings.Contains(buf, " ")
	// check handlers
	for _, h := range s {
		if hasSpace {
			candidates := h.Complete(buf)
			if candidates != nil {
				return candidates
			}
			continue
		}
		// complete the command itself
		if util.HasPrefixFold(h.Name, buf) {
			c = append(c, h.Name)
		}
		for _, a := range h.Aliases {
			if util.HasPrefixFold(a, buf) {
				c = append(c, a)
			}
		}
	}

	return c
}

func (s Set) RunHelp(arg string) error {
	if arg == "" {
		for _, h := range s {
			fmt.Println(h.String())
		}
	} else {
		h := s.Match(arg)
		if h == nil {
			fmt.Printf("%s is not a known command or alias.\nRun `help` for a list of available commands.\n", arg)
			return nil
		}
		fmt.Println(h.GoString())
	}
	return nil
}
