package control

import (
	"fmt"
	"strings"

	"github.com/insomnimus/libman/control/store"
	"github.com/insomnimus/libman/handler"
	"github.com/zmb3/spotify"
)

func defaultImportHandlers() handler.Set {
	const (
		Help = iota
		Create
		Replace
		Play
	)

	hand := func(c uint8, name, help, usage string, aliases ...string) handler.Handler {
		return handler.Handler{
			Cmd:      c,
			Name:     name,
			Usage:    usage,
			Help:     help,
			About:    help,
			Aliases:  aliases,
			Complete: completeNothing,
			Run:      func(string) error { panic("internal error: .Run() called on a dummy handler") },
		}
	}

	set := handler.Set{
		hand(
			Create,
			"create",
			"Create a new playlist, with the contents of the imported one.",
			"create [name]",
			"save",
		),
		hand(
			Replace,
			"replace",
			"Replace an existing playlist with the imported one.",
			"replace [name]",
			"repl",
		),
		hand(
			Play,
			"play",
			"Play tracks from the imported playlist.",
			"play",
			"pl",
		),
		hand(
			Help,
			"help",
			"List available commands.",
			"help [command]",
		),
	}

	set.Find(Help).Complete = set.CompleteHelp
	set.Find(Help).Run = set.RunHelp
	updateCache()
	set.Find(Replace).Complete = newWordCompleter(cache.ownedNames()...)

	return set
}

func handleExportPlaylist(arg string) error {
	p := choosePlaylist(arg)
	if p == nil {
		return nil
	}
	filename := readString("Enter the full path where the file will be saved (.json extension will be appended): ")
	if filename == "" {
		fmt.Println("cancelled")
		return nil
	}

	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	if err := p.makeFull(); err != nil {
		return err
	}

	err := store.ExportTo(p.FullPlaylist, filename)
	if err == nil {
		fmt.Printf("Exported %s to %s.\n", p.Name, filename)
	}
	return err
}

func handleImportPlaylist(arg string) error {
	const (
		Help = iota
		Create
		Replace
		Play
	)

	if arg == "" {
		arg = readString("Enter the path of a previously exported playlist: ")
		if arg == "" {
			fmt.Println("cancelled")
			return nil
		}
	}

	p, err := store.ImportFrom(arg)
	if err != nil {
		return err
	}
	fmt.Printf("Imported %s.\n", p.Name)

	if len(p.Tracks) == 0 {
		fmt.Printf("%s has no tracks in it, no action can be taken.\n", p.Name)
		return nil
	}

	fmt.Println("Type `help` for a list of available actions.")

	set := defaultImportHandlers()
	set.Find(Play).Run = func(string) error {
		uris := make([]spotify.URI, len(p.Tracks))
		for i, t := range p.Tracks {
			uris[i] = t.Track.URI
		}
		err := client.PlayOpt(&spotify.PlayOptions{
			URIs: uris,
		})
		if err == nil {
			fmt.Printf("Playing tracks from the imported playlist %s.\n", p.Name)
			isPlaying = true
		}
		return err
	}

	set.Find(Replace).Run = func(arg string) error {
		if arg == "" {
			fmt.Println("Choose a playlist to replace.")
		}
		pl := choosePlaylist(arg)
		if pl == nil {
			return nil
		}

		if readBool("Are you sure you want to replace all tracks in %s with those in %s?", pl.Name, p.Name) {
			ids := make([]spotify.ID, len(p.Tracks))
			for i, t := range p.Tracks {
				ids[i] = t.Track.ID
			}
			err := client.ReplacePlaylistTracks(pl.ID, ids...)
			if err == nil {
				fmt.Printf("Replaced every track in %s.\n", pl.Name)
				pl.Tracks.Tracks = p.Tracks
			}
			return err
		} else {
			fmt.Println("cancelled")
			return nil
		}
	}

	set.Find(Create).Run = func(arg string) error {
		if len(p.Tracks) == 0 {
			fmt.Println("The imported playlist has no tracks. If you still want to create a playlist, use the `create` command.")
			return nil
		}
		if arg == "" {
			arg = readString("Playlist name (%s): ", p.Name)
			if arg == "" {
				arg = p.Name
			}
		}
		desc := readString("Playlist description: ")
		pub := readBool("Should the playlist be public?")
		if !readBool("Create new playlist %s and import tracks?", arg) {
			fmt.Println("cancelled")
			return nil
		}

		pl, err := client.CreatePlaylistForUser(user.ID, arg, desc, pub)
		if err != nil {
			return err
		}
		fmt.Printf("Created new playlist %q.\n", pl.Name)
		ids := make([]spotify.ID, len(p.Tracks))
		for i, t := range p.Tracks {
			ids[i] = t.Track.ID
		}
		_, err = client.AddTracksToPlaylist(pl.ID, ids...)
		if err != nil {
			fmt.Println("Failed to add new tracks to the playlist.")
			return err
		}
		pl.Tracks.Tracks = p.Tracks
		cache.insertFull(0, *pl)
		fmt.Printf("Successfully imported all the tracks to %s.\n", pl.Name)
		return nil
	}

	for {
		rl.SetCompleter(set.Complete)
		reply, cancelled := readPrompt(false, "imported %s$ ", p.Name)
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if reply == "" {
			continue
		}
		command, arg := splitCmd(reply)
		h := set.Match(command)
		if h == nil {
			fmt.Printf("%s is not a known command or alias.\nRun `help` for a list of available commands.\n", command)
			continue
		}
		return h.Run(arg)
	}
}
