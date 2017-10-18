package main

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	IncludeExtensions  []string `toml:"include_extensions"`
	ExcludeDirectories []string `toml:"exclude_directories"`
	Name               string   `toml:"name"`

	includeExtensionMap map[string]bool
	excludeDirectoryMap map[string]bool
}

var DefaultConfig = Config{
	IncludeExtensions: []string{
		"css",
		"go",
		"java",
		"js",
		"scss",
		"ts",
	},
	ExcludeDirectories: []string{
		".git",
		"node_modules",
	},
}

func GetConfig() *Config {
	config := &Config{}
	config.Merge(&DefaultConfig, readConfig())
	return config
}

func (c *Config) Merge(configs ...*Config) {
	for _, config := range configs {
		if config != nil {
			if config.ExcludeDirectories != nil {
				c.ExcludeDirectories = config.ExcludeDirectories
				c.excludeDirectoryMap = nil
			}
			if config.IncludeExtensions != nil {
				c.IncludeExtensions = config.IncludeExtensions
				c.includeExtensionMap = nil
			}
			if config.Name != "" {
				c.Name = config.Name
			}
		}
	}
}

// Augment the default configuration with
func readConfig() *Config {
	dir, err := homedir.Dir()
	if err != nil {
		return nil
	}
	configPath := path.Join(dir, ".todorc")
	if _, err = os.Stat(configPath); err != nil {
		return nil
	}
	var config Config
	if _, err = toml.DecodeFile(configPath, &config); err != nil {
		return nil
	}
	return &config
}

func (c *Config) ShouldScanDir(directory string) bool {
	if c.excludeDirectoryMap == nil {
		c.excludeDirectoryMap = make(map[string]bool)
		for _, dir := range c.ExcludeDirectories {
			c.excludeDirectoryMap[dir] = true
		}
	}
	_, shouldExcludeDirectory := c.excludeDirectoryMap[directory]
	return !shouldExcludeDirectory
}

func (c *Config) ShouldScanFile(filename string) bool {
	if c.includeExtensionMap == nil {
		c.includeExtensionMap = make(map[string]bool)
		for _, ext := range c.IncludeExtensions {
			c.includeExtensionMap[ext] = true
		}
	}
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return false
	}
	extension := parts[len(parts)-1]
	_, shouldIncludeFile := c.includeExtensionMap[extension]
	return shouldIncludeFile
}

func (c *Config) Validate() (bool, []error) {
	var errs []error
	if c.Name == "" {
		errs = append(errs, errors.New("Name to search for is required"))
	}
	return len(errs) == 0, errs
}
