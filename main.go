package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/insomnimus/libman/auth"
	"github.com/insomnimus/libman/config"
	"github.com/insomnimus/libman/control"
	// "github.com/vrischmann/userdir"
)

const VERSION = "0.5.3"

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	path := os.Getenv("LIBMAN_CONFIG_PATH")
	if path == "" {
		cfg, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("error: could not determine user config dir: %s\n", err)
		}
		path = filepath.Join(cfg, "libman.toml")
	}

	c, err := config.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	creds, err := auth.Login(c)
	if err != nil {
		log.Fatal(err)
	}
	// load the commands from the rc file, if it exists
	var commands []string
	if c.RCFile != "" {
		if _, err := os.Stat(c.RCFile); err == nil {
			data, err := os.ReadFile(c.RCFile)
			if err != nil {
				log.Fatalf("error reading from the libman startup script file at %s:\n%s\n", c.RCFile, err)
			}
			commands = strings.Split(string(data), "\n")
		}
	}

	// The liner package traps sigint
	// but not while it's being blocked, so trap it here as well.
	signal.Notify(make(chan os.Signal), os.Interrupt)
	go control.Start(creds.Client, creds.User, c.Prompt, commands)
	<-control.Terminator
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
		err = os.WriteFile(c.CacheFile, data, 0o600)
		if err != nil {
			log.Fatalf("error saving the token to the cache file: %s", err)
		}
	} else {
		log.Println("warning: the access token is not saved because no cache file is specified")
	}
}
