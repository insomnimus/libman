package control

import (
	"fmt"
	"strconv"

	"github.com/insomnimus/libman/handler/cmd"
	"github.com/zmb3/spotify"
)

func handleCreatePlaylist(arg string) error {
	rl.SetCompleter(func(string) []string { return nil })
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

	if !readBool("confirm: creating new playlist %s, proceed?", arg) {
		fmt.Println("cancelled")
		return nil
	}

	pl, err := client.CreatePlaylistForUser(user.ID, arg, desc, pub)
	if err != nil {
		return err
	}

	fmt.Printf("created new playlist %q\n", arg)
	cache.insertFull(0, *pl)
	return nil
}

func handleDeletePlaylist(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}
	msg := "Are you sure you want to delete %s?"
	if pl.isFollowed {
		msg = "Are you sure you want to unfollow %s?"
	}
	if !readBool(msg, pl.Name) {
		fmt.Println("cancelled")
		return nil
	}

	err := client.UnfollowPlaylist(spotify.ID(pl.Owner.ID), pl.ID)
	if err != nil {
		return err
	}

	msg = "Deleted %s.\n"
	if pl.isFollowed {
		msg = "Unfollowed %s.\n"
	}
	fmt.Printf(msg, pl.Name)
	cache.remove(pl.ID)
	if lastPl != nil && pl.ID == lastPl.ID {
		lastPl = nil
	}
	return nil
}

func handleEditPlaylistDetails(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}
	return pl.editDetails()
}

func handleEditPlaylist(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}
	return pl.editTracks()
}

func updateCache() error {
	if cache == nil {
		// page, err := client.CurrentUsersPlaylists()

		cache = make(PlaylistCache, 0)
		var pls []spotify.SimplePlaylist

		for {
			limit := 50
			page, err := client.GetPlaylistsForUserOpt(user.ID, &spotify.Options{
				Limit: &limit,
			})
			if err != nil {
				return err
			}
			pls = append(pls, page.Playlists...)
			if len(page.Playlists) < limit {
				break
			}
		}
		for _, p := range pls {
			if p.Owner.ID != string(user.ID) {
				cache.pushFollowed(p)
			} else {
				cache.pushSimple(p)
			}
		}
	}
	return nil
}

// func updateAlbumCache() error {
// if savedAlbums == nil {
// page, err := client.CurrentUsersAlbums()
// if err != nil {
// return err
// }
// savedAlbums = make(AlbumCache, len(page.Albums))
// for i, a := range page.Albums {
// savedAlbums[i] = a.FullAlbum.SimpleAlbum
// }
// }
// return nil
// }

func likeTrack(t *spotify.FullTrack) error {
	if err := updateLibraryCache(); err != nil {
		return err
	}
	if libraryCache.contains(t.ID) {
		fmt.Printf("%s is already in your library.\n", t.Name)
		return nil
	}
	err := client.AddTracksToLibrary(t.ID)
	if err == nil {
		fmt.Printf("Saved %s to the library.\n", t.Name)
		libraryCache.push(*t)
	}
	return err
}

func followArtist(a *spotify.FullArtist) error {
	follows, err := client.CurrentUserFollows("artist", a.ID)
	if err != nil {
		return err
	}
	if len(follows) > 0 && follows[0] {
		fmt.Printf("You are already following %s.\n", a.Name)
		return nil
	}
	err = client.FollowArtist(a.ID)
	if err == nil {
		fmt.Printf("Followed %s.\n", a.Name)
	}
	return err
}

func saveAlbum(spotify.SimpleAlbum) error {
	/******** There's currently no function to do this in zmb3/spotify
	if err := updateAlbumCache(); err != nil{
		return err
	}
	if savedAlbums.contains(a.ID) {
		fmt.Printf("%s is already in your library.\n", a.Name)
	}
	********/
	fmt.Println("Not yet implemented.")
	return nil
}

func followPlaylist(p *Playlist) error {
	err := client.FollowPlaylist(spotify.ID(p.Owner.ID), p.ID, true)
	if err == nil {
		fmt.Printf("Followed %s.\n", p.Name)
	}
	return err
}

func updateLibraryCache() error {
	if libraryCache == nil {
		limit := 50
		page, err := client.CurrentUsersTracksOpt(&spotify.Options{
			Limit: &limit,
		})
		if err != nil {
			return err
		}
		libraryCache = make(LibraryCache, 0, len(page.Tracks))
		if len(page.Tracks) == 0 {
			return nil
		}
		if len(page.Tracks) < 50 {
			for _, t := range page.Tracks {
				libraryCache.push(t.FullTrack)
			}
			return nil
		}
		// There might be more than 50 tracks, page until api returns none.
		for offset := 1; len(page.Tracks) > 0; offset++ {
			ofs := offset * limit
			page, err = client.CurrentUsersTracksOpt(&spotify.Options{
				Limit:  &limit,
				Offset: &ofs,
			})
			if err != nil {
				return err
			}
			for _, t := range page.Tracks {
				libraryCache.push(t.FullTrack)
			}
		}
	}
	return nil
}

func handleLikePlaying(string) error {
	t, err := getPlaying()
	if err != nil {
		return err
	}
	if err := updateLibraryCache(); err != nil {
		return err
	}
	if libraryCache.contains(t.ID) {
		fmt.Printf("%s is already in your library.\n", t.Name)
		return nil
	}
	err = client.AddTracksToLibrary(t.ID)
	if err == nil {
		libraryCache.push(*t)
		fmt.Printf("Saved %s to your library.\n", t.Name)
	}
	return err
}

func handleDislikePlaying(string) error {
	t, err := getPlaying()
	if err != nil {
		return err
	}
	if err := updateLibraryCache(); err != nil {
		return err
	}
	if !libraryCache.contains(t.ID) {
		fmt.Printf("%s is not in your library.\n", t.Name)
		return nil
	}
	err = client.RemoveTracksFromLibrary(t.ID)
	if err == nil {
		libraryCache.removeByID(t.ID)
		fmt.Printf("Removed %s from your library.\n", t.Name)
	}
	return err
}

func handlePlayTopTracks(arg string) error {
	limit := 20
	if arg != "" {
		n, err := strconv.Atoi(arg)
		if err != nil || n <= 0 {
			handlers.ShowUsage(cmd.PlayTopTracks)
			return nil
		}
		if n > 50 {
			fmt.Println("Error: the max limit is 50.")
			return nil
		}
		limit = n
	}

	page, err := client.CurrentUsersTopTracksOpt(&spotify.Options{Limit: &limit})
	if err != nil {
		return err
	}
	if len(page.Tracks) == 0 {
		fmt.Println("Spotify says you have no top tracks...")
		return nil
	}

	uris := make([]spotify.URI, len(page.Tracks))
	fmt.Printf("Playing your top %d tracks (recent):\n", len(page.Tracks))
	for i, t := range page.Tracks {
		fmt.Printf("-\t%s by %s\n", t.Name, joinArtists(t.Artists))
		uris[i] = t.URI
	}

	err = client.PlayOpt(&spotify.PlayOptions{
		URIs: uris,
	})
	if err == nil {
		isPlaying = true
	}
	return err
}

func choosePlaylist(arg string) *Playlist {
	// if the cache is nil, initialize it
	if err := updateCache(); err != nil {
		fmt.Printf("Error updating playlist cache: %s.\n", err)
		return nil
	}

	if len(cache) == 0 {
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

	totalPages := len(cache) / 20
	if totalPages == 0 || len(cache)%20 != 0 {
		totalPages++
	}

	min := func(a, b int) int {
		if a < b {
			return a
		} else {
			return b
		}
	}

	display := func(page int) {
		start := page * 20
		if start > len(cache) {
			return
		}
		end := min(page*20+20, len(cache))

		fmt.Printf("Displaying page %d of %d.\n", page+1, totalPages)
		for i, p := range cache[start:end] {
			fmt.Printf("%3d | %s\n", i+start, p.Name)
		}
	}
	page := 0

	nextPage := func() {
		if (page+1)*20 >= len(cache) {
			fmt.Println("Already on the last page.")
			return
		}
		page++
		display(page)
	}

	prevPage := func() {
		if page == 0 {
			fmt.Println("Already on the first page.")
			return
		}
		page--
		display(page)
	}

	display(0)

	for {
		input, cancelled := readPrompt(false, "[%d-%d, < and > for paging, ^c to cancel]: ", page*20, min(page*20+20, len(cache)))
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

		if input == ">" {
			nextPage()
			continue
		}
		if input == "<" {
			prevPage()
			continue
		}

		n, err := strconv.Atoi(input)
		if err != nil ||
			n < 20*page ||
			n >= min(20*page+20, len(cache)) {
			fmt.Printf("Please enter a value between %d and %d, use `<` and `>` for paging and press ctrl+c to cancel.\n",
				page*20, min(page*20+20, len(cache)))
			continue
		}
		return &cache[n]
	}
}
