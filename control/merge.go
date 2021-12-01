package control

import (
	"fmt"

	"github.com/zmb3/spotify"
)

func mergeTracks(pls ...*Playlist) ([]spotify.PlaylistTrack, []spotify.ID, error) {
	var tracks []spotify.PlaylistTrack
	var ids []spotify.ID
	added := make(map[spotify.ID]struct{})

	for _, pl := range pls {
		if err := pl.makeFull(); err != nil {
			return nil, nil, err
		}
		for _, t := range pl.Tracks.Tracks {
			if _, contains := added[t.Track.ID]; !contains {
				added[t.Track.ID] = struct{}{}
				tracks = append(tracks, t)
				ids = append(ids, t.Track.ID)
			}
		}
	}

	return tracks, ids, nil
}

func choosePlaylists(args string) []*Playlist {
	if err := updateCache(); err != nil {
		fmt.Printf("Error updating playlist cache: %s.\n", err)
		return nil
	}
	if len(cache) == 0 {
		fmt.Println("you don't seem to have any playlists")
		return nil
	}

	var pls []*Playlist
	fmt.Println("Add at least 2 playlists into the list.\nWhen you're done, press enter without any input.")
	if args != "" {
		pl := cache.findByName(args)
		if pl == nil {
			fmt.Printf("You don't seem to have a playlist named %q.\n", args)
			return nil
		}
		pls = append(pls, pl)
	}

LOOP:
	for {
		rl.SetCompleter(newCompleter(cache.names()...))
		input, cancelled := readPrompt(false, "playlist name: ")
		if cancelled {
			fmt.Println("cancelled")
			return nil
		}
		if input == "" {
			if len(pls) > 1 {
				return pls
			}
			fmt.Println("You need to add at least 2 playlists in order to merge. You can press ctrl-c to abort.")
			continue
		}
		pl := cache.findByName(input)
		if pl == nil {
			fmt.Println("You don't have a playlist by that name.")
			continue
		}
		for _, added := range pls {
			if added.ID == pl.ID {
				fmt.Printf("You already added %s.", pl.Name)
				continue LOOP
			}
		}
		pls = append(pls, pl)
		fmt.Printf("Added %s.\n", pl.Name)
	}
}

func handleMerge(args string) error {
	pls := choosePlaylists(args)
	if pls == nil {
		return nil
	}
	rl.SetCompleter(func(string) []string { return nil })
	fmt.Printf("Merging %d playlists.\n", len(pls))
	var name string
	for {
		input, cancelled := readPrompt(false, "enter new playlist name: ")
		if (cancelled || input == "") && readBool("Are you sure you want to abort?") {
			fmt.Println("cancelled")
			return nil
		}
		name = input
		break
	}
	desc := readString("Playlist description: ")
	pub := readBool("Should the playlist be public?")
	tracks, ids, err := mergeTracks(pls...)
	if err != nil {
		return err
	}
	pl, err := client.CreatePlaylistForUser(user.ID, name, desc, pub)
	if err != nil {
		return err
	}
	for i := 0; i < len(ids); i += 100 {
		upto := i + 100
		if upto > len(ids) {
			upto = len(ids)
		}
		_, err = client.AddTracksToPlaylist(pl.ID, ids[i:upto]...)
		if err != nil {
			fmt.Println("Failed to populate the new playlist.")
			return err
		}
	}
	pl.Tracks.Tracks = tracks
	cache.insertFull(0, *pl)
	fmt.Printf("Success. %s is ready.\n", name)
	return nil
}
