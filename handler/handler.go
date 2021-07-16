package handler

import (
	"fmt"
	"libman/handler/cmd"
	"strings"
)

type Handler struct {
	Name    string
	Aliases []string
	Cmd     cmd.Cmd
	Help    string
	About   string
	Usage   string
	Run     func(string) error
}

type Set []Handler

func (s Set) Find(c cmd.Cmd) *Handler {
	for _, h := range s {
		if h.Cmd == c {
			return &h
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

func (s Set) ShowUsage(c cmd.Cmd) {
	if h := s.Find(c); h != nil {
		fmt.Printf("usage:\n  %s\n", h.Usage)
	}
}
