package control

import (
	"fmt"
	"github.com/insomnimus/libman/handler/cmd"
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

	n := readNumber(0, len(tracks))
	if n < 0 {
		fmt.Println("cancelled")
		return nil
	}

	return playTrack(&tracks[n])
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

	n := readNumber(0, len(albs))
	if n < 0 {
		fmt.Println("cancelled")
		return nil
	}

	return playAlbum(&albs[n])
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

	n := readNumber(0, len(arts))
	if n < 0 {
		fmt.Println("cancelled")
		return nil
	}

	return playArtist(&arts[n])
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
		fmt.Printf("#%3d | %s from %s\n", i, p.Name, p.Owner.DisplayName)
	}

	n := readNumber(0, len(pls))
	if n < 0 {
		fmt.Println("cancelled")
		return nil
	}

	return playPlaylist(&pls[n])
}
