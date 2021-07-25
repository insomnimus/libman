package control

import (
	"fmt"

	"github.com/insomnimus/libman/handler/cmd"
)

func handlePFTrack(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstTrack)
		return nil
	}

	tracks, err := searchTrack(arg)
	if err != nil {
		return err
	}
	if len(tracks) == 0 {
		fmt.Printf("no result for %q\n", arg)
		return nil
	}
	// append to tracks history
	Hist.AppendTrack(tracks[0].Name)
	return playTrack(&tracks[0])
}

func handlePFAlbum(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstAlbum)
		return nil
	}
	albs, err := searchAlbum(arg)
	if err != nil {
		return err
	}
	if len(albs) == 0 {
		fmt.Printf("no result for %q\n", arg)
		return nil
	}

	Hist.AppendAlbum(albs[0].Name)

	return playAlbum(&albs[0])
}

func handlePFArtist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstArtist)
		return nil
	}

	arts, err := searchArtist(arg)
	if err != nil {
		return err
	}

	if len(arts) == 0 {
		fmt.Printf("no result for %s\n", arg)
		return nil
	}

	Hist.AppendArtist(arts[0].Name)

	return playArtist(&arts[0])
}

func handlePFPlaylist(arg string) error {
	if arg == "" {
		handlers.ShowUsage(cmd.PlayFirstPlaylist)
		return nil
	}

	pls, err := searchPlaylist(arg)
	if err != nil {
		return err
	}

	if len(pls) == 0 {
		fmt.Printf("no result for %q", arg)
		return nil
	}

	Hist.AppendPlaylist(pls[0].Name)

	return playPlaylist(&pls[0])
}
