package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	ID          string `toml:"client_id" comment:"The $LIBMAN_ID will override this field if set."`
	Secret      string `toml:"client_secret" comment:"The $LIBMAN_SECRET will override this field if set."`
	RedirectURI string `toml:"redirect_uri" comment:"The $LIBMAN_REDIRECT_URI will override this field if set."`
	CacheFile   string `toml:"cache_file_path" comment:"Full file path of the libman token cache file.\nThe $LIBMAN_CACHE_PATH will override this, if set."`
	RCFile      string `toml:"libmanrc_path" comment:"The location of the startup script file. $LIBMAN_RC_PATH overrides this field if set."`
	HistFile    string `toml:"history_file" comment:"File where artist, album and track history will be saved to. $LIBMAN_HIST_FILE overrides this, if set."`

	Prompt string `toml:"prompt" comment:"The libman shell prompt." commented:"true" default:"@libman>"`
}

func DefaultConfig() Config {
	return Config{
		ID:          os.Getenv("LIBMAN_ID"),
		Secret:      os.Getenv("LIBMAN_SECRET"),
		RedirectURI: os.Getenv("LIBMAN_REDIRECT_URI"),
		CacheFile:   os.Getenv("LIBMAN_CACHE_PATH"),
		RCFile:      RCPath(),
		Prompt:      "@libman>",
	}
}

func Load(path string) (*Config, error) {
	// create file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c := DefaultConfig()
		data, err := toml.Marshal(c)
		if err != nil {
			return nil, err
		}
		data = append([]byte("# libman configuration file\n"), data...)

		return &c, os.WriteFile(path, data, 0o600)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = toml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("malformed config file: %w", err)
	}

	// id := os.Getenv("LIBMAN_ID")
	// secret := os.Getenv("LIBMAN_SECRET")
	// uri := os.Getenv("LIBMAN_REDIRECT_URI")
	// cache := os.Getenv("LIBMAN_CACHE_PATH")

	// if id != "" {
	// c.ID = id
	// }
	// if secret != "" {
	// c.Secret = secret
	// }
	// if uri != "" {
	// c.RedirectURI = uri
	// }
	// if cache != "" {
	// c.CacheFile = cache
	// }

	return &c, nil
}

func RCPath() string {
	if s := os.Getenv("LIBMAN_RC_PATH"); s != "" {
		return s
	}
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(u.HomeDir, ".libmanrc")
}
