package control

import (
	"fmt"
	"github.com/insomnimus/libman/alias"
	"github.com/insomnimus/libman/handler"
	"github.com/zmb3/spotify"
	"regexp"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

var (
	client *spotify.Client
	user   *spotify.PrivateUser
	device *spotify.PlayerDevice

	prompt       = "@libman>"
	userAliases  = new(alias.Set)
	handlers     handler.Set
	cache        *PlaylistCache
	lastPl       *Playlist
	isPlaying    bool
	shuffleState bool
	repeatState  = "off"

	reVol = regexp.MustCompile(`^\s*(\-|\+)\s*([0-9]+)\s*$`)

	rl         *liner.State
	Terminator = make(chan bool, 1)
)

func SetHandlers(h handler.Set) {
	handlers = h
}

func Start(
	c *spotify.Client,
	u *spotify.PrivateUser,
	p string,
	commands []string,
) {
	rl = liner.NewLiner()
	rl.SetCtrlCAborts(true)
	defer rl.Close()
	client = c
	user = u
	if prompt != "" {
		prompt = p
	}

	// execute the startup commands
	if len(commands) > 0 {
		var err error
		for _, cmd := range commands {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" || strings.HasPrefix(cmd, "#") {
				continue
			}
			cmd = expandAlias(cmd)
			cmd, arg := splitCmd(cmd)
			if h := handlers.Match(cmd); h != nil {
				err = h.Run(arg)
			} else {
				err = fmt.Errorf("%s is not a known command or alias", cmd)
			}
			if err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
		}
	}

	var input string
	var err error
	var cancelled bool
	for {
		rl.SetCompleter(completeCommand)
		input, cancelled = readPrompt(prompt + " ")
		if cancelled {
			continue
		}
		input = expandAlias(input)
		if input == "" {
			err = togglePlay()
		} else if m := reVol.FindStringSubmatch(input); len(m) == 3 {
			n, _ := strconv.Atoi(m[2])
			if m[1] == "-" {
				n = -n
			}
			err = adjustVolume(n)
		} else {
			cmd, arg := splitCmd(input)
			h := handlers.Match(cmd)
			if h == nil {
				fmt.Printf("%s is not a known command or alias.\nRun `help` for a list of available commands.\n", cmd)
				continue
			}
			err = h.Run(arg)
		}
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}
}
