package control

import (
	"fmt"
	"strings"

	"github.com/insomnimus/libman/handler"
	"github.com/insomnimus/libman/handler/scmd"
)

func defaultSTrackHandlers() handler.Set {
	hand := func(
		c uint8,
		name string,
		about string,
		usage string,
		help string,
		aliases ...string,
	) handler.Handler {
		return handler.Handler{
			Cmd:     c,
			Name:    name,
			About:   about,
			Usage:   usage,
			Help:    help,
			Aliases: aliases,
			Run: func(string) error {
				panic("called handler.Run which was nil")
			},
			Complete: completeNothing,
		}
	}

	set := handler.Set{
		hand(
			scmd.Play,
			"play",
			"Play a track.",
			"play <N>",
			"Play a track from the list. You can also just enter its number.",
			"pl",
		),
		hand(
			scmd.Save,
			"save",
			"Save a track to a playlist.",
			"save <N> <playlist>",
			"Save a track from the list to one of your playlists.",
			"add",
		),
		hand(
			scmd.Like,
			"like",
			"Add a track to your favourites folder.",
			"like <N>",
			"Add a track from the list to your favorites folder.",
			"fav", "fave",
		),
		hand(
			scmd.Queue,
			"queue",
			"Add a track to your playing queue.",
			"queue <N>",
			"Add a track from the list to your playing queue.",
			"q", "que",
		),
		hand(
			scmd.Help,
			"help",
			"Display help about a command.",
			"help [command]",
			"Display help about a command or list available commands.",
		),
	}

	set.Find(scmd.Help).Run = set.RunHelp
	set.Find(scmd.Help).Complete = newWordCompleter(set.CommandsAndAliases(), "help")
	set.Find(scmd.Save).Complete = func(buf string) []string {
		buf = strings.TrimPrefix(buf, " ")
		spaceLast := strings.HasSuffix(buf, " ")
		command, arg := splitCmd(buf)
		if !strings.EqualFold(command, "save") && !strings.EqualFold(command, "add") {
			return nil
		}
		// do not complete if the <N> field is not there
		if !spaceLast && !strings.Contains(arg, " ") {
			return nil
		}
		_, arg = splitCmd(arg)
		// complete arg (playlist name)
		if err := updateCache(); err != nil {
			return nil
		}
		pls := make([]string, 0, len(cache))
		if arg == "" {
			for _, p := range cache {
				pls = append(pls, buf+p.Name)
			}
			return pls
		}

		for _, p := range cache {
			if hasPrefixFold(p.Name, arg) {
				pls = append(pls, fmt.Sprintf("%s %s", buf, p.Name))
			}
		}
		// return nil if there are no candidates
		if len(pls) == 0 {
			return nil
		}
		return pls
	}

	return set
}

func defaultSArtistHandlers() handler.Set {
	hand := func(
		c uint8,
		name string,
		about string,
		usage string,
		help string,
		aliases ...string,
	) handler.Handler {
		return handler.Handler{
			Cmd:     c,
			Name:    name,
			About:   about,
			Usage:   usage,
			Help:    help,
			Aliases: aliases,
			Run: func(string) error {
				panic("called handler.Run which was nil")
			},
			Complete: completeNothing,
		}
	}

	set := handler.Set{
		hand(
			scmd.Play,
			"play",
			"Play an artist.",
			"play <N>",
			"Play an artist from the list.",
			"pl",
		),
		hand(
			scmd.Follow,
			"follow",
			"Follow an artist.",
			"follow <N>",
			"Follow an artist from the list.",
		),
		hand(
			scmd.Help,
			"help",
			"Display help about a command.",
			"help [command]",
			"Display help about a command or list available commands.",
		),
	}

	set.Find(scmd.Help).Complete = newWordCompleter(set.CommandsAndAliases(), "help")
	set.Find(scmd.Help).Run = set.RunHelp

	return set
}

func defaultSAlbumHandlers() handler.Set {
	hand := func(
		c uint8,
		name string,
		about string,
		usage string,
		help string,
		aliases ...string,
	) handler.Handler {
		return handler.Handler{
			Cmd:     c,
			Name:    name,
			About:   about,
			Usage:   usage,
			Help:    help,
			Aliases: aliases,
			Run: func(string) error {
				panic("called handler.Run which was nil")
			},
			Complete: completeNothing,
		}
	}

	set := handler.Set{
		hand(
			scmd.Play,
			"play",
			"Play an album.",
			"play <N>",
			"Play an album from the list.",
			"pl",
		),
		hand(
			scmd.Save,
			"save",
			"Save an album to your library.",
			"save <N>",
			"Save an album from the list to your library.",
		),
		hand(
			scmd.Queue,
			"queue",
			"Add an albums tracks to the playing queue.",
			"queue <N>",
			"Add all the tracks of an album from the list to your playing queue.",
			"q", "que",
		),
		hand(
			scmd.Help,
			"help",
			"Display help about a command.",
			"help [command]",
			"Display help about a command or list available commands.",
		),
	}

	set.Find(scmd.Help).Run = set.RunHelp
	set.Find(scmd.Help).Complete = newWordCompleter(set.CommandsAndAliases(), "help")

	return set
}

func defaultSPlaylistHandlers() handler.Set {
	hand := func(
		c uint8,
		name string,
		about string,
		usage string,
		help string,
		aliases ...string,
	) handler.Handler {
		return handler.Handler{
			Cmd:     c,
			Name:    name,
			About:   about,
			Usage:   usage,
			Help:    help,
			Aliases: aliases,
			Run: func(string) error {
				panic("called handler.Run which was nil")
			},
			Complete: completeNothing,
		}
	}

	set := handler.Set{
		hand(
			scmd.Play,
			"play",
			"Play a playlist.",
			"play <N>",
			"Play a playlist from the list.",
			"pl",
		),
		hand(
			scmd.Follow,
			"follow",
			"Follow a playlist.",
			"follow <N>",
			"Follow a playlist from the list.",
		),
		hand(
			scmd.Help,
			"help",
			"display help about a command.",
			"help [command]",
			"Display help about a command or list available commands.",
		),
	}

	set.Find(scmd.Help).Run = set.RunHelp
	set.Find(scmd.Help).Complete = newWordCompleter(set.CommandsAndAliases(), "help")

	return set
}
