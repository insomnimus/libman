package control

import (
	"fmt"
	"github.com/zmb3/spotify"
	"libman/handler"
	"regexp"
	"strconv"
)

var (
	client       *spotify.Client
	user         *spotify.PrivateUser
	device       *spotify.PlayerDevice
	prompt       = "@libman>"
	handlers     handler.Set
	cache        *PlaylistCache
	lastPl       *Playlist
	isPlaying    bool
	shuffleState bool
	repeatState  = "off"

	reVol = regexp.MustCompile(`^\s*(\-|\+)\s*([0-9]+)\s*$`)
)

func SetHandlers(h handler.Set) {
	handlers = h
}

func Start(c *spotify.Client, u *spotify.PrivateUser, p string) {
	client = c
	user = u
	prompt = p
	var input string
	var err error
	for {
		input = readString(prompt + " ")
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
