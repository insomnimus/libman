package control

import (
	"fmt"
	"strconv"

	"github.com/insomnimus/libman/handler/cmd"
	"github.com/insomnimus/libman/handler/scmd"
	"github.com/zmb3/spotify"
)

func handleSTrack(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.SearchTrack)
		return nil
	}
	tracks, err := searchTrack(arg)
	if err != nil {
		return err
	}
	if len(tracks) == 0 {
		fmt.Printf("No result for %s.\n", arg)
		return nil
	}

	// print the results
	for i, t := range tracks {
		fmt.Printf("%-2d | %s by %s\n", i, t.Name, joinArtists(t.Artists))
	}

	return sTrackPage(tracks)
}

func handleSAlbum(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.SearchAlbum)
		return nil
	}
	albs, err := searchAlbum(arg)
	if err != nil {
		return err
	}

	if len(albs) == 0 {
		fmt.Printf("No result for %s.\n", arg)
		return nil
	}

	for i, a := range albs {
		fmt.Printf("#%2d | %s by %s\n", i, a.Name, joinArtists(a.Artists))
	}

	return sAlbumPage(albs)
}

func handleSArtist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.SearchArtist)
		return nil
	}
	arts, err := searchArtist(arg)
	if err != nil {
		return err
	}

	if len(arts) == 0 {
		fmt.Printf("No result for %s.\n", arg)
		return nil
	}
	for i, a := range arts {
		fmt.Printf("#%2d | %s\n", i, a.Name)
	}

	return sArtistPage(arts)
}

func handleSPlaylist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.SearchPlaylist)
		return nil
	}
	pls, err := searchPlaylist(arg)
	if err != nil {
		return err
	}

	if len(pls) == 0 {
		fmt.Printf("No result for %s.\n", arg)
		return nil
	}
	for i, p := range pls {
		fmt.Printf("#%2d | %s from %s\n", i, p.Name, p.Owner.DisplayName)
	}
	return sPlaylistPage(pls)
}

func sTrackPage(tracks []spotify.FullTrack) error {
	fmt.Println("Type `help` for a list of available commands.")

	for {
		rl.SetCompleter(sTrackHandlers.Complete)
		input, cancelled := readPrompt(false, "command: ")
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if input == "" {
			if err := togglePlay(); err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
			continue
		}
		command, arg := splitCmd(input)
		h := sTrackHandlers.Match(command)
		if h == nil {
			if arg == "" {
				n, err := strconv.Atoi(command)
				if err != nil {
					fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", arg)
					continue
				}
				if n < 0 || n >= len(tracks) {
					fmt.Printf("Please enter a value between 0 and %d.\n", len(tracks))
					continue
				}
				return playTrack(&tracks[n])
			}
			fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", arg)
			continue
		}
		switch h.Cmd {
		case scmd.Help:
			h.Run(arg)
		case scmd.Play:
			n, ok := parseNumber(arg, len(tracks), h.Usage)
			if ok {
				return playTrack(&tracks[n])
			}
		case scmd.Like:
			n, ok := parseNumber(arg, len(tracks), h.Usage)
			if ok {
				return likeTrack(&tracks[n])
			}
		case scmd.Queue:
			n, ok := parseNumber(arg, len(tracks), h.Usage)
			if ok {
				return queueTrack(&tracks[n])
			}
		case scmd.Save:
			if arg == "" {
				fmt.Println(h.Usage)
				break
			}
			num, name := splitCmd(arg)
			if name == "" {
				fmt.Println(h.Usage)
				break
			}
			n, ok := parseNumber(num, len(tracks), h.Usage)
			if ok {
				pl := choosePlaylist(name)
				if pl != nil {
					return pl.addTrack(tracks[n])
				}
			}
		default:
			panic(fmt.Sprintf("internal error: unhandled case %s\n", h.Name))
		}
	}
}

func sArtistPage(artists []spotify.FullArtist) error {
	fmt.Println("Type `help` for a list of available commands.")
	for {
		rl.SetCompleter(sArtistHandlers.Complete)
		input, cancelled := readPrompt(false, "command: ")
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if input == "" {
			if err := togglePlay(); err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
			continue
		}

		command, arg := splitCmd(input)
		h := sArtistHandlers.Match(command)
		if h == nil {
			if arg == "" {
				n, err := strconv.Atoi(command)
				if err != nil {
					fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
					continue
				}
				if n < 0 || n >= len(artists) {
					fmt.Printf("Please enter a value between 0 and %d.\n", len(artists))
					continue
				}
				return playArtist(&artists[n])
			}
			fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
			continue
		}

		if h.Cmd == scmd.Help {
			h.Run(arg)
			continue
		}
		n, ok := parseNumber(arg, len(artists), h.Usage)
		if !ok {
			continue
		}
		switch h.Cmd {
		case scmd.Play:
			return playArtist(&artists[n])
		case scmd.Follow:
			return followArtist(&artists[n])
		default:
			panic(fmt.Sprintf("internal error: handler missing case: %s\n", h.Name))
		}
	}
}

func sAlbumPage(albums []spotify.SimpleAlbum) error {
	fmt.Println("Run `help` for a list of available commands.")
	for {
		rl.SetCompleter(sAlbumHandlers.Complete)
		input, cancelled := readPrompt(false, "command: ")
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if input == "" {
			if err := togglePlay(); err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
			continue
		}

		command, arg := splitCmd(input)
		h := sAlbumHandlers.Match(command)
		if h == nil {
			if arg == "" {
				n, err := strconv.Atoi(command)
				if err != nil {
					fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
					continue
				}
				if n < 0 || n >= len(albums) {
					fmt.Printf("Please enter a value between 0 and %d.\n", len(albums))
					continue
				}
				return playAlbum(&albums[n])
			}
			fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
			continue
		}

		if h.Cmd == scmd.Help {
			h.Run(arg)
			continue
		}

		n, ok := parseNumber(arg, len(albums), h.Usage)
		if !ok {
			continue
		}
		switch h.Cmd {
		case scmd.Play:
			return playAlbum(&albums[n])
		case scmd.Save:
			return saveAlbum(albums[n])
		case scmd.Queue:
			return queueAlbum(&albums[n])
		default:
			panic(fmt.Sprintf("internal error: unhandled handler case: %s\n", h.Name))
		}
	}
}

func sPlaylistPage(pls []Playlist) error {
	fmt.Println("Run `help` for a list of available commands.")
	for {
		rl.SetCompleter(sPlaylistHandlers.Complete)
		input, cancelled := readPrompt(false, "command: ")
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if input == "" {
			if err := togglePlay(); err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
			continue
		}

		command, arg := splitCmd(input)
		h := sPlaylistHandlers.Match(command)
		if h == nil {
			if arg == "" {
				n, err := strconv.Atoi(command)
				if err != nil {
					fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
					continue
				}
				if n < 0 || n >= len(pls) {
					fmt.Printf("Please enter a value between 0 and %d.\n", len(pls))
					continue
				}
				return playPlaylist(&pls[n])
			}
			fmt.Printf("%s is not a known command, alias or a number.\nRun `help` for a list of available commands.\n", command)
			continue
		}

		if h.Cmd == scmd.Help {
			h.Run(arg)
			continue
		}

		n, ok := parseNumber(arg, len(pls), h.Usage)
		if !ok {
			continue
		}
		switch h.Cmd {
		case scmd.Play:
			return playPlaylist(&pls[n])
		case scmd.Follow:
			return followPlaylist(&pls[n])
		default:
			panic(fmt.Sprintf("internal error: unhandled handler case: %s\n", h.Name))
		}
	}
}
