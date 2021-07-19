package control

import (
	"fmt"
	"github.com/zmb3/spotify"
	"strconv"
)

func handleCreatePlaylist(arg string) error {
	rl.SetCompleter(completeNothing)
	fmt.Println("creating new playlist")
	if arg == "" {
		arg = readString("playlist name: ")
		if arg == "" {
			fmt.Println("cancelled")
			return nil
		}
	} else {
		fmt.Printf("playlist name: %s\n", arg)
	}

	desc := readString("playlist description: ")
	pub := readBool("should the playlist be public?")

	if !readBool("confirm\ncreating new playlist %s, proceed?", arg) {
		fmt.Println("cancelled")
		return nil
	}

	pl, err := client.CreatePlaylistForUser(user.ID, arg, desc, pub)
	if err != nil {
		return err
	}

	fmt.Printf("created new playlist %q", arg)
	cache.insertFull(0, *pl)
	return nil
}

func handleDeletePlaylist(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}

	if !readBool("are you sure you want to delete %s?", pl.Name) {
		fmt.Println("cancelled")
		return nil
	}

	err := client.UnfollowPlaylist(spotify.ID(pl.Owner.ID), pl.ID)
	if err != nil {
		return err
	}

	fmt.Printf("deleted %s\n", pl.Name)
	cache.remove(pl.ID)
	if lastPl != nil && pl.ID == lastPl.ID {
		lastPl = nil
	}
	return nil
}

func handleEditPlaylist(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}
	return pl.editDetails()
}

func updateCache() error {
	if cache == nil {
		page, err := client.CurrentUsersPlaylists()
		if err != nil {
			return err
		}
		cache = new(PlaylistCache)
		for _, p := range page.Playlists {
			cache.pushSimple(p)
		}
	}
	return nil
}

func choosePlaylist(arg string) *Playlist {
	// if the cache is nil, initialize it
	err := updateCache()
	if err != nil {
		fmt.Printf("Error updating playlist cache: %s.\n", err)
		return nil
	}

	if len(*cache) == 0 {
		fmt.Println("you don't seem to have any playlists")
		return nil
	}

	if arg != "" {
		pl := cache.findByName(arg)
		if pl == nil {
			fmt.Printf("You don't seem to have a playlist named %q.\n", arg)
		}
		return pl
	}

	if len(*cache) > PlaylistPageSize {
		return _pagePlaylist()
	}

	for i, p := range *cache {
		fmt.Printf("%-2d | %s\n", i, p.Name)
	}

	n := readNumber(0, len(*cache))
	if n == -1 {
		fmt.Println("cancelled")
		return nil
	}
	return cache.get(n)
}

func _pagePlaylist() *Playlist {
	var page int
	nextPage := func() {
		page++
		startFrom := page * PlaylistPageSize
		if startFrom < len(*cache) {
			fmt.Println("Already on the last page.")
			page--
			return
		}
		upTo := (page + 1) * PlaylistPageSize
		if upTo > len(*cache) {
			upTo = len(*cache)
		}
		for i, p := range (*cache)[startFrom:upTo] {
			fmt.Printf("#%3d | %s\n", startFrom+i, p.Name)
		}
	}

	prevPage := func() {
		if page == 0 {
			fmt.Println("Already on the first page.")
			return
		}
		page--
		startFrom := page * PlaylistPageSize
		upTo := startFrom + 1
		if upTo > len(*cache) {
			upTo = len(*cache)
		}

		for i, p := range (*cache)[startFrom:upTo] {
			fmt.Printf("#%3d | %s\n", i+startFrom, p.Name)
		}
	}
	page = -1
	nextPage()
	fmt.Println("Type `<` and `>` for the prev/next pages.")

	for {
		startFrom := page * PlaylistPageSize
		upTo := startFrom + PlaylistPageSize
		if upTo > len(*cache) {
			upTo = len(*cache)
		}
		text := readString("[%d-%d, blank to cancel]: ", startFrom, upTo)
		if text == ">" {
			nextPage()
		} else if text == "<" {
			prevPage()
		} else if text == "" {
			fmt.Println("cancelled")
			return nil
		} else {
			n, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("Please enter a value between %d and %d, blank to cancel and <> for paging.\n", startFrom, upTo)
				continue
			}
			if n < startFrom || n >= upTo {
				fmt.Printf("Please enter a value between %d and %d.\n", startFrom, upTo)
				continue
			}
			return cache.get(n)
		}
	}
}
