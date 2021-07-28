package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var configFields = map[string]func(*Config) interface{}{
	"cache-file":   func(c *Config) interface{} { return &c.CacheFile },
	"id":           func(c *Config) interface{} { return &c.ID },
	"secret":       func(c *Config) interface{} { return &c.Secret },
	"redirect-uri": func(c *Config) interface{} { return &c.RedirectURI },
	"history-file": func(c *Config) interface{} { return &c.HistFile },
	"history-size": func(c *Config) interface{} { return &c.HistSize },
	"rc-file":      func(c *Config) interface{} { return &c.RCFile },
	"prompt":       func(c *Config) interface{} { return &c.Prompt },
	"data-home":    func(c *Config) interface{} { return &c.DataHome },
}

func runHelp(c *cli.Context) {
	fmt.Printf(`libman %s: Configure libman.
Usage:
  libman %s <FIELD> [VALUE]
or
  libman %s [OPTIONS]

OPTIONS:
  -l, --list: List available config fields.
  -h, --help: Show this message and exit.

To set a field, specify it's value (VALUE).
To view a field, omit the VALUE
`, c.Command.Name, c.Command.Name, c.Command.Name)
	os.Exit(0)
}

func Command() *cli.Command {
	return &cli.Command{
		Name:     "config",
		Usage:    "Set libman configuration fields.",
		HideHelp: true,
		Action:   run,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List available config fields.",
			},
			&cli.BoolFlag{
				Name:    "help",
				Aliases: []string{"h"},
				Usage:   "Show this message and exist.",
			},
		},
	}
}

func run(c *cli.Context) error {
	if c.Bool("help") {
		runHelp(c)
		os.Exit(0)
	}

	if err := runConfig(c); err != nil {
		log.Fatalf("error: %s\n", err)
	}
	os.Exit(0)
	return nil
}

func runConfig(ctx *cli.Context) error {
	if ctx.Bool("list") {
		path, err := ConfigPath()
		if err != nil {
			for k := range configFields {
				fmt.Println(k)
			}
			return nil
		}
		c, err := Load(path)
		if err != nil {
			for k := range configFields {
				fmt.Println(k)
			}
			return nil
		}
		for k, fn := range configFields {
			switch val := fn(c).(type) {
			case *string:
				fmt.Printf("%s = %q\n", k, *val)
			case *int:
				fmt.Printf("%s = %d\n", k, *val)
			default:
				panic("internal error: unhandled type switch case in config/cmd.go")
			}
		}
		return nil
	}

	args := ctx.Args().Slice()
	if len(args) == 0 {
		runHelp(ctx)
		return nil
	}

	if len(args) > 2 {
		return fmt.Errorf("too many arguments, run with --help for the usage")
	}
	key := strings.ToLower(args[0])
	fn, ok := configFields[key]
	if !ok {
		return fmt.Errorf("%s is not a valid config field, run with --list for a list of available fields", key)
	}
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	c, err := Load(path)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		switch val := fn(c).(type) {
		case *string:
			fmt.Printf("%s = %s\n", key, *val)
		case *int:
			fmt.Printf("%s = %d\n", key, *val)
		default:
			panic("internal error: unhandled type switch in config/cmd.go")
		}
		return nil
	}
	var msg string
	switch val := fn(c).(type) {
	case *string:
		*val = args[1]
		msg = fmt.Sprintf("%s = %q", key, args[1])
	case *int:
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("%s is not a valid value for %s (expected integer)", args[1], key)
		}
		*val = n
		msg = fmt.Sprintf("%s = %d", key, n)
	default:
		panic("internal error: unhandled type switch in config/cmd.go")
	}

	err = c.Save(path)
	if err == nil {
		fmt.Println(msg)
	}
	return err
}
