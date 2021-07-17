package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
)

type Config struct {
	ID          string `toml:"client_id" comment:"The LIBMAN_ID env variable will override this field if set."`
	Secret      string `toml:"client_secret" comment:"The LIBMAN_SECRET env variable will override this field if set."`
	RedirectURI string `toml:"redirect_uri" comment:"The LIBMAN_REDIRECT_URI env variable will override this field if set."`
	CacheFile   string `toml:"cache_file_path" comment:"Full file path of the libman token cache file.\nThe LIBMAN_CACHE_PATH env variable will override this, if set."`
	Prompt      string `toml:"prompt" comment:"The libman shell prompt." commented:"true" default:"@libman>"`
}

func DefaultConfig() Config {
	return Config{
		ID:          os.Getenv("LIBMAN_ID"),
		Secret:      os.Getenv("LIBMAN_SECRET"),
		RedirectURI: os.Getenv("LIBMAN_REDIRECT_URI"),
		CacheFile:   os.Getenv("LIBMAN_CACHE_PATH"),
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

		return &c, os.WriteFile(path, data, 0600)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = toml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("malformed config file: %e\n", err)
	}

	id := os.Getenv("LIBMAN_ID")
	secret := os.Getenv("LIBMAN_SECRET")
	uri := os.Getenv("LIBMAN_REDIRECT_URI")
	cache := os.Getenv("LIBMAN_CACHE_PATH")

	if id != "" {
		c.ID = id
	}
	if secret != "" {
		c.Secret = secret
	}
	if uri != "" {
		c.RedirectURI = uri
	}
	if cache != "" {
		c.CacheFile = cache
	}

	return &c, nil
}
