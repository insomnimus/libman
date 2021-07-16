package control

import (
	"fmt"
	"github.com/zmb3/spotify"
	"libman/handler/cmd"
	"strings"
)

func togglePlayback() error {
	if err := updateDevice(); err != nil {
		return err
	}
	isPlaying = !isPlaying
	var err error
	if isPlaying {
		err = client.PlayOpt(&spotify.PlayOptions{
			DeviceID: &device.ID,
		})
	} else {
		err = client.PauseOpt(&spotify.PlayOptions{
			DeviceID: &device.ID,
		})
	}
	return err
}

func setVolume(n int) error {
	if err := updateDevice(); err != nil {
		return err
	}
	if n < 0 {
		n = 0
	} else if n > 100 {
		n = 100
	}
	err := client.VolumeOpt(n, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	device.Volume = n
	return nil
}

func adjustVolume(n int) error {
	if err := updateDevice(); err != nil {
		return err
	}
	n += device.Volume
	if n > 100 {
		n = 100
	} else if n < 0 {
		n = 0
	}
	err := client.VolumeOpt(n, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	device.Volume = n
	return nil
}

func getActiveDevice() (*spotify.PlayerDevice, error) {
	devs, err := client.PlayerDevices()
	if err != nil {
		return nil, err
	}
	for _, d := range devs {
		if d.Active && !d.Restricted {
			return &d, nil
		}
	}

	return nil, nil
}

func updateDevice() error {
	if device == nil {
		dev, err := getActiveDevice()
		if err != nil {
			return err
		}
		device = dev
	}
	return nil
}

func playNext() error {
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.NextOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	isPlaying = true
	return nil
}

func playPrev() error {
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PreviousOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	isPlaying = true
	return nil
}

func handleShuffle(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	var b bool
	switch strings.ToLower(arg) {
	case "on", "true", "yes":
		b = true
	case "off", "no", "false":
		b = false
	case "":
		b = !shuffleState
	default:
		handlers.ShowUsage(cmd.Shuffle)
		return nil
	}
	err := client.ShuffleOpt(b, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	if b {
		fmt.Println("shuffle on")
	} else {
		fmt.Println("shuffle off")
	}
	shuffleState = b
	return nil
}

func handleRepeat(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	var r string
	switch strings.ToLower(arg) {
	case "off", "false", "no":
		r = "off"
	case "on", "track", "song":
		r = "track"
	case "context":
		r = "context"
	default:
		handlers.ShowUsage(cmd.Repeat)
		return nil
	}

	err := client.RepeatOpt(r, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("repeat = %s\n", r)
	repeatState = r
	return nil
}
