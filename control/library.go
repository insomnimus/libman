package control

import (
	"fmt"
	"github.com/zmb3/spotify"
)

func createPlaylist(arg string) error {
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

func deletePlaylist(arg string) error {
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
	return nil
}

func choosePlaylist(arg string) *Playlist {
	// if the cache is nil, initialize it
	if cache == nil {
		page, err := client.CurrentUsersPlaylists()
		if err != nil {
			fmt.Println("error fetching user playlists: ", err)
			return nil
		}
		cache = new(PlaylistCache)
		for _, p := range page.Playlists {
			cache.pushSimple(p)
		}
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

	for i, p := range *cache {
		fmt.Printf("%.2d | %s\n", i, p.Name)
	}

	n := readNumber(0, len(*cache))
	if n == -1 {
		fmt.Println("cancelled")
		return nil
	}
	return cache.get(n)
}
