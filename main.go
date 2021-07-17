package main

import (
	"encoding/json"
	"libman/auth"
	"libman/config"
	"libman/control"
	"log"
	"os"
	"os/signal"

	"github.com/vrischmann/userdir"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	path := os.Getenv("LIBMAN_CONFIG_PATH")
	if path == "" {
		path = userdir.GetDataHome()
	}

	c, err := config.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	creds, err := auth.Login(c)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	control.SetHandlers(control.DefaultHandlers())

	go control.Start(creds.Client, creds.User, c.Prompt)
	<-ch

	// on windows, control + c makes the player switch pause state, do it again
	if os.PathSeparator == '\\' {
		control.TogglePlay()
	}

	// save the token if there's a cache file specified
	if c.CacheFile != "" {
		token, err := creds.Client.Token()
		if err != nil {
			log.Fatalf("error retreiving token: %s", err)
		}
		data, err := json.MarshalIndent(token, "", "\t")
		if err != nil {
			log.Fatalf("error serializing token as json: %s", err)
		}
		err = os.WriteFile(c.CacheFile, data, 0600)
		if err != nil {
			log.Fatalf("error saving the token to the cache file: %s", err)
		}
	} else {
		log.Println("warning: the access token is not saved because no cache file is specified")
	}
}
