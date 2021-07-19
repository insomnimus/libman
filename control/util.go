package control

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/insomnimus/libman/util"

	"github.com/zmb3/spotify"
)

// aliases for compatibility
var (
	hasPrefixFold = util.HasPrefixFold
	splitCmd      = util.SplitCmd
)

func readBool(format string, args ...interface{}) bool {
	rl.SetCompleter(completeBool)
	for {
		reply := readString(format+" [y/n]: ", args...)
		switch strings.ToLower(reply) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("please enter yes or no")
		}
	}
}

func readPrompt(addToHistory bool, format string, args ...interface{}) (reply string, cancelled bool) {
	var err error
	reply, err = rl.Prompt(fmt.Sprintf(format, args...))
	if errors.Is(err, io.EOF) {
		Terminator <- true
		select {}
	}
	if err != nil {
		return "", true
	}
	if addToHistory {
		rl.AppendHistory(reply)
	}
	return strings.TrimSpace(reply), false
}

func readString(format string, args ...interface{}) string {
	reply, _ := readPrompt(false, format, args...)
	return reply
}

func trackQuery(s string) string {
	if strings.Contains(s, "::") {
		split := strings.SplitN(s, "::", 2)
		track := strings.TrimSpace(split[0])
		if len(split) != 2 {
			return "track:" + track
		}
		artist := strings.TrimSpace(split[1])
		return fmt.Sprintf("track:%s artist:%s", track, artist)
	}

	if strings.Contains(s, " by ") {
		split := strings.SplitN(s, " by ", 2)
		if len(split) != 2 {
			return s
		}
		track := strings.TrimSpace(split[0])
		artist := strings.TrimSpace(split[1])
		return fmt.Sprintf("track:%s artist:%s", track, artist)
	}

	return s
}

func albumQuery(s string) string {
	if strings.Contains(s, "::") {
		split := strings.SplitN(s, "::", 2)
		alb := strings.TrimSpace(split[0])
		if len(split) != 2 {
			return "album:" + alb
		}
		artist := strings.TrimSpace(split[1])
		return fmt.Sprintf("album:%s artist:%s", alb, artist)
	}

	if strings.Contains(s, " by ") {
		split := strings.SplitN(s, " by ", 2)
		if len(split) != 2 {
			return s
		}
		alb := strings.TrimSpace(split[0])
		artist := strings.TrimSpace(split[1])
		return fmt.Sprintf("album:%s artist:%s", alb, artist)
	}

	return s
}

func readNumber(min, max int) int {
	for {
		reply := readString("[%d-%d, -1 or blank to cancel]: ", min, max)
		if reply == "" || reply == "-1" {
			return -1
		}

		n, err := strconv.Atoi(reply)
		if err != nil {
			fmt.Println("invalid input, please enter again")
			continue
		}

		if n < min || n >= max {
			fmt.Printf("please enter a number between %d and %d\n", min, max)
			continue
		}

		return n
	}
}

func joinArtists(arts []spotify.SimpleArtist) string {
	switch len(arts) {
	case 0:
		return ""
	case 1:
		return arts[0].Name
	case 2:
		return fmt.Sprintf("%s and %s", arts[0].Name, arts[1].Name)
	default:
		names := make([]string, len(arts)-1)
		for i, a := range arts[:len(arts)-1] {
			names[i] = a.Name
		}
		return fmt.Sprintf("%s and %s", strings.Join(names, ", "), arts[len(arts)-1].Name)
	}
}

func expandAlias(s string) string {
	als := userAliases.Inner()
	if als != nil {
		lower := strings.ToLower(s)
		for _, a := range als {
			left := strings.ToLower(a.Left)
			if strings.HasPrefix(lower, left) {
				// must be separated by a space
				if len(left) == len(s) {
					return a.Right
				}
				if s[len(left)] == ' ' {
					return fmt.Sprintf("%s %s", a.Right, s[len(left):])
				}
			}
		}
	}
	return s
}
