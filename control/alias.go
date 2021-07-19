package control

import (
	"fmt"
	"strings"

	"github.com/insomnimus/libman/handler/cmd"
)

func handleAlias(arg string) error {
	if arg == "" {
		if userAliases.Len() == 0 {
			fmt.Println("You don't have any aliases.")
			return nil
		}

		for _, a := range userAliases.Sorted() {
			fmt.Println(a.String())
		}
		return nil
	}
	if !strings.Contains(arg, "=") {
		a, ok := userAliases.Get(arg)
		if ok {
			fmt.Println(a.String())
		} else {
			fmt.Printf("No alias %q found.\n", arg)
		}
		return nil
	}

	split := strings.SplitN(arg, "=", 2)
	left := strings.TrimSpace(split[0])
	right := strings.TrimSpace(split[1])

	switch {
	case left == "" && right == "":
		handlers.ShowUsage(cmd.Alias)
	case right == "":
		ok := userAliases.Unset(left)
		if ok {
			fmt.Printf("Unset alias %s.\n", left)
		} else {
			fmt.Printf("%s is not a known user set alias.\n", left)
		}
	case left == "":
		handlers.ShowUsage(cmd.Alias)
	default:
		userAliases.Set(left, right)
	}
	return nil
}
