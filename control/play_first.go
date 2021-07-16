package control

import (
	"fmt"
	"github.com/zmb3/spotify"
	"libman/handler/cmd"
)

func handlePFTrack(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstTrack)
		return nil
	}

	q := trackQuery(arg)
	page, err := client.Search(q, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}

	if page.Tracks == nil || len(page.Tracks.Tracks) == 0 {
		fmt.Printf("no result for %q\n", q)
		return nil
	}

	track := &page.Tracks.Tracks[0]
	return playTrack(track)
}

func handlePlayFAlbum(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstAlbum)
		return nil
	}

	q := albumQuery(arg)
	page, err := client.Search(q, spotify.SearchTypeAlbum)
	if err != nil {
		return err
	}
	if page.Albums == nil || len(page.Albums.Albums) == 0 {
		fmt.Printf("no result for %q\n", q)
		return nil
	}

	alb := &page.Albums.Albums[0]
	return playAlbum(alb)
}

func handlePlayFArtist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstArtist)
		return nil
	}

	page, err := client.Search(arg, spotify.SearchTypeArtist)
	if err != nil {
		return err
	}

	if page.Artists == nil || len(page.Artists.Artists) == 0 {
		fmt.Printf("no result for %s\n", arg)
		return nil
	}

	art := &page.Artists.Artists[0]

	return playArtist(art)
}

func handlePlayFPlaylist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstPlaylist)
		return nil
	}

	page, err := client.Search(arg, spotify.SearchTypePlaylist)
	if err != nil {
		return err
	}

	if page.Playlists == nil || len(page.Playlists.Playlists) == 0 {
		fmt.Printf("no result for %q", arg)
		return nil
	}

	pl := plFromSimple(page.Playlists.Playlists[0])
	return playPlaylist(&pl)
}
