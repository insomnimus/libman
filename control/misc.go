package control

import (
	"fmt"

	"github.com/insomnimus/libman/handler/cmd"
)

func handlePrompt(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.Prompt)
		return nil
	}
	prompt = arg
	return nil
}

func handleHelp(arg string) error {
	if arg == "" {
		for _, h := range handlers {
			fmt.Println(h.String())
		}
		fmt.Println("\nYou can also paste a spotify link to play from.")
	} else {
		h := handlers.Match(arg)
		if h == nil {
			fmt.Printf("%s is not a known command or alias.\nRun `help` for a list of available commands.\n", arg)
			return nil
		}
		fmt.Println(h.GoString())
	}
	return nil
}
