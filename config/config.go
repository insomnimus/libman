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
	DataHome    string `toml:"data_home" comment:"$LIBMAN_DATA_HOME will override this, if set. The place where the exported playlist will be saved to, if the path given is not absolute."`

	HistSize int    `toml:"history_size" comment:"History size, applies to artist/album/track/playlist history independently." default:"66"`
	Prompt   string `toml:"prompt" comment:"The libman shell prompt." default:"@libman>"`

	ConfigPath string `toml:"-"`
}

func DefaultConfig() Config {
	dataHome := os.Getenv("LIBMAN_DATA_HOME")
	if dataHome == "" {
		p, err := os.UserCacheDir()
		if err == nil {
			dataHome = filepath.Join(p, "libman")
			os.MkdirAll(dataHome, 0o700)
		}
	}

	configPath := os.Getenv("LIBMAN_CONFIG_PATH")
	if configPath == "" {
		p, err := os.UserConfigDir()
		if err == nil {
			p = filepath.Join(p, "libman")
			os.MkdirAll(p, 0o644)
			configPath = filepath.Join(p, "libman.toml")
		}
	}

	return Config{
		ID:          os.Getenv("LIBMAN_ID"),
		Secret:      os.Getenv("LIBMAN_SECRET"),
		RedirectURI: os.Getenv("LIBMAN_REDIRECT_URI"),
		CacheFile:   os.Getenv("LIBMAN_CACHE_PATH"),
		RCFile:      RCPath(),
		Prompt:      "@libman>",
		HistSize:    66,
		ConfigPath:  configPath,
		DataHome:    dataHome,
	}
}

func Load(path string) (*Config, error) {
	// create file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c := DefaultConfig()
		c.ConfigPath = path
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

	c.ConfigPath = path

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

func ConfigPath() (string, error) {
	path := os.Getenv("LIBMAN_CONFIG_PATH")
	if path != "" {
		return path, nil
	}
	path, err := os.UserConfigDir()
	if err == nil {
		path = filepath.Join(path, "libman.toml")
	}
	return path, err
}

func (c *Config) Save(path string) error {
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
