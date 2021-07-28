package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/insomnimus/libman/config"
	"github.com/urfave/cli/v2"
)

func printFlags(flags []cli.Flag) {
	for _, f := range flags {
		switch f := f.(type) {
		case *cli.IntFlag:
			fmt.Printf("  %s: %s",
				joinFlags(f.Name, f.Aliases),
				f.Usage)

			if len(f.EnvVars) > 0 {
				fmt.Printf("    [$%s]", f.EnvVars[0])
			}
			fmt.Println()
		case *cli.StringFlag:
			fmt.Printf("  %s: %s",
				joinFlags(f.Name, f.Aliases),
				f.Usage)

			if len(f.EnvVars) > 0 {
				fmt.Printf("    [$%s]", f.EnvVars[0])
			}
			fmt.Println()
		case *cli.BoolFlag:
			fmt.Printf("%s: %s\n", joinFlags(f.Name, f.Aliases),
				f.Usage)

		default:
			panic("internal error: unhandled case switch")
		}
	}
}

func joinFlags(name string, aliases []string) string {
	var shorts []string
	var longs []string
	if len(name) == 1 {
		shorts = append(shorts, "-"+name)
	} else {
		longs = append(longs, "--"+name)
	}
	for _, f := range aliases {
		if len(f) == 1 {
			shorts = append(shorts, "-"+f)
		} else {
			longs = append(longs, "--"+f)
		}
	}

	if len(longs) == 0 {
		return strings.Join(shorts, ", ")
	} else if len(shorts) == 0 {
		return strings.Join(longs, ", ")
	} else {
		return strings.Join(append(shorts, longs...), ", ")
	}
}

// handleHelp is called when the help command or flag is set for the main command.
// Because quite frankly, urfave/cli's auto-generated help sucks.
func handleHelp(c *cli.Context) error {
	fmt.Printf("%s: %s\n", c.App.Name, c.App.Usage)
	fmt.Printf("Usage:\n  %s [OPTIONS]\n", c.App.Name)
	if len(c.App.VisibleCommands()) > 0 {
		fmt.Printf("or\n  %s <SUBCOMMAND>\n", c.App.Name)
	}

	fmt.Println("OPTIONS:")

	printFlags(c.App.VisibleFlags())

	commands := c.App.VisibleCommands()
	if len(commands) > 0 {
		fmt.Println("Subcommands:")
		for _, cmd := range commands {
			fmt.Printf("  %s: %s\n", cmd.Name, cmd.Usage)
		}
	}
	os.Exit(0)
	return nil
}

func handleVersion(c *cli.Context) error {
	fmt.Printf("%s version %s\n", c.App.Name, VERSION)
	os.Exit(0)
	return nil
}

func run(c *cli.Context) error {
	if c.Bool("help") {
		handleHelp(c)
	}
	if c.Bool("version") {
		handleVersion(c)
	}
	return nil
}

func configFromArgs() (*config.Config, error) {
	var (
		id         string
		cache      string
		secret     string
		configPath string
		rc         string
		hist       string
		redirect   string
		prompt     string
		histSize   *int
		dataHome   string
	)

	app := &cli.App{
		Name:        "libman",
		Usage:       "An interactive spotify shell.",
		Version:     VERSION,
		Description: "Libman is an interactive spotify shell.\nIt lets you control your spotify playback, manage your library and more.",
		HideHelp:    true,
		HideVersion: true,
		Action:      run,
		Commands:    []*cli.Command{config.Command()},
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
				Name:        "data-home",
				Destination: &dataHome,
				Usage:       "The path to the libman data directory where the exported playlists will be stored.",
				EnvVars:     []string{"LIBMAN_DATA_HOME"},
			},
			&cli.StringFlag{
				Name:        "hist-file",
				Usage:       "The file where your search history will be saved to, for autocompletion.",
				Aliases:     []string{"t"},
				TakesFile:   true,
				EnvVars:     []string{"LIBMAN_HIST_FILE"},
				Destination: &hist,
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
				Destination: &prompt,
			},
			&cli.BoolFlag{
				Name:    "help",
				Aliases: []string{"h"},
				Usage:   "Show this message and exit.",
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "Show the installed version and exit.",
			},
			&cli.IntFlag{
				Name:        "hist-size",
				Usage:       "History size, recommended to keep under 100.",
				Destination: histSize,
				EnvVars:     []string{"LIBMAN_HIST_SIZE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return nil, err
	}

	if configPath == "" {
		p, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("could not locate a config file: %w", err)
		}
		configPath = filepath.Join(p, "libman.toml")
	}

	// if all the flags are set, no need to read the config file
	if id != "" &&
		secret != "" &&
		cache != "" &&
		redirect != "" &&
		rc != "" &&
		hist != "" &&
		histSize != nil &&
		prompt != "" &&
		dataHome != "" {
		return &config.Config{
			ID:          id,
			Secret:      secret,
			RedirectURI: redirect,
			CacheFile:   cache,
			RCFile:      rc,
			HistFile:    hist,
			HistSize:    *histSize,
			Prompt:      prompt,
			ConfigPath:  configPath,
			DataHome:    dataHome,
		}, nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if rc == "" {
		rc = config.RCPath()
	}

	if id != "" {
		cfg.ID = id
	}
	if secret != "" {
		cfg.Secret = secret
	}
	if redirect != "" {
		cfg.RedirectURI = redirect
	}
	if cache != "" {
		cfg.CacheFile = cache
	}
	if hist != "" {
		cfg.HistFile = hist
	}
	if rc != "" {
		cfg.RCFile = rc
	}
	if prompt != "" {
		cfg.Prompt = prompt
	}
	if dataHome != "" {
		cfg.DataHome = dataHome
		if cfg.CacheFile == "" {
			cfg.CacheFile = filepath.Join(dataHome, "libman_token_cachejson")
		}
	}
	return cfg, nil
}
