package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func FromArgs(version string) (*Config, error) {
	var (
		id            string
		cache         string
		secret        string
		configPath    string
		rc            string
		redirect      string
		prompt        string
		isMainCommand bool
	)

	app := &cli.App{
		Name:                  "libman",
		Usage:                 "An interactive spotify shell.",
		Version:               version,
		Description:           "Libman is an interactive spotify shell.\nIt lets you control your spotify playback, manage your library and more.",
		CustomAppHelpTemplate: HelpTemplate,
		Action:                func(*cli.Context) error { isMainCommand = true; return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "id",
				Destination: &id,
				DefaultText: "Read from env or config file",
				Usage:       "The spotify client id.",
				Aliases:     []string{"i"},
				EnvVars:     []string{"LIBMAN_ID"},
			},
			&cli.StringFlag{
				Name:        "secret",
				DefaultText: "Read from env or config file",
				Destination: &secret,
				Usage:       "The spotify client secret.",
				Aliases:     []string{"s"},
				EnvVars:     []string{"LIBMAN_SECRET"},
			},
			&cli.StringFlag{
				Name:        "redirect-uri",
				DefaultText: "Read from env or config file",
				Destination: &redirect,
				Usage:       "The redirect uri, must be localhost.",
				Aliases:     []string{"u"},
				EnvVars:     []string{"LIBMAN_REDIRECT_URI"},
			},
			&cli.StringFlag{
				Name:        "config-file",
				DefaultText: "Read from env or config file",
				Destination: &configPath,
				Usage:       "A file where unspecified fields will be read from.",
				Aliases:     []string{"c"},
				EnvVars:     []string{"LIBMAN_CONFIG_PATH"},
				TakesFile:   true,
			},
			&cli.StringFlag{
				Name:        "cache-file",
				DefaultText: "Read from env or config file",
				Destination: &cache,
				Usage:       "A full path to a file where the session token will be saved/ read from.",
				Aliases:     []string{"C"},
				EnvVars:     []string{"LIBMAN_CACHE_PATH"},
				TakesFile:   true,
			},
			&cli.StringFlag{
				Name:        "rc-file",
				DefaultText: "Read from env or config file",
				Destination: &rc,
				Usage:       "A text file containing libman commands to be ran after startup.",
				Aliases:     []string{"r"},
				EnvVars:     []string{"LIBMAN_RC_PATH"},
				TakesFile:   true,
			},
			&cli.StringFlag{
				Name:        "prompt",
				Usage:       "The libman shell prompt, can be set in-app.",
				Value:       "@libman>",
				Destination: &prompt,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return nil, err
	}
	if !isMainCommand {
		os.Exit(0)
	}
	// if all the flags are set, no need to read the config file
	if id != "" &&
		secret != "" &&
		cache != "" &&
		redirect != "" &&
		rc != "" &&
		prompt != "" {
		return &Config{
			ID:          id,
			Secret:      secret,
			RedirectURI: redirect,
			CacheFile:   cache,
			RCFile:      rc,
			Prompt:      prompt,
		}, nil
	}

	if configPath == "" {
		p, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("could not locate a config file: %w", err)
		}
		configPath = filepath.Join(p, "libman.toml")
	}

	config, err := Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if rc == "" {
		rc = rcPath()
	}

	if id != "" {
		config.ID = id
	}
	if secret != "" {
		config.Secret = secret
	}
	if redirect != "" {
		config.RedirectURI = redirect
	}
	if cache != "" {
		config.CacheFile = cache
	}
	if rc != "" {
		config.RCFile = rc
	}
	if prompt != "" {
		config.Prompt = prompt
	}

	return config, nil
}
