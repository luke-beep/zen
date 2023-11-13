package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

var Config config

type filterList struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Url     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

type config struct {
	Filter struct {
		FilterLists []filterList `json:"filterLists"`
	} `json:"filter"`
	Certmanager struct {
		CAInstalled bool `json:"caInstalled"`
	} `json:"certmanager"`
	Proxy struct {
		Port int `json:"port"`
	} `json:"proxy"`
	ConfigDir string `json:"-"`
	DataDir   string `json:"-"`
}

func (c *config) Save() error {
	configData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	configFile := path.Join(c.ConfigDir, "config.json")
	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetFilterLists returns the list of enabled filter lists.
// Used on the frontend to display the list of filter lists.
func (c *config) GetFilterLists() []filterList {
	return c.Filter.FilterLists
}

// AddFilterList adds a new filter list to the list of enabled filter lists.
// Used on the frontend to add a new filter list.
func (c *config) AddFilterList(list filterList) string {
	c.Filter.FilterLists = append(c.Filter.FilterLists, list)
	if err := c.Save(); err != nil {
		fmt.Printf("failed to save config: %v", err)
		return err.Error()
	}
	return ""
}

// RemoveFilterList removes a filter list from the list of enabled filter lists.
// Used on the frontend to remove a filter list.
func (c *config) RemoveFilterList(url string) string {
	for i, filterList := range c.Filter.FilterLists {
		if filterList.Url == url {
			c.Filter.FilterLists = append(c.Filter.FilterLists[:i], c.Filter.FilterLists[i+1:]...)
			break
		}
	}
	if err := c.Save(); err != nil {
		fmt.Printf("failed to save config: %v", err)
		return err.Error()
	}
	return ""
}

// ToggleFilterList toggles the enabled state of a filter list.
// Used on the frontend to toggle the enabled state of a filter list.
func (c *config) ToggleFilterList(url string, enabled bool) string {
	for i, filterList := range c.Filter.FilterLists {
		if filterList.Url == url {
			c.Filter.FilterLists[i].Enabled = enabled
			break
		}
	}
	if err := c.Save(); err != nil {
		fmt.Printf("failed to save config: %v", err)
		return err.Error()
	}
	return ""
}

// GetPort returns the port the proxy is set to listen on.
// Used on the frontend in the settings manager.
func (c *config) GetPort() int {
	return c.Proxy.Port
}

// SetPort sets the port the proxy is set to listen on.
// Used on the frontend in the settings manager.
func (c *config) SetPort(port int) string {
	c.Proxy.Port = port
	if err := c.Save(); err != nil {
		fmt.Printf("failed to save config: %v", err)
		return err.Error()
	}
	return ""
}

//go:embed default-config.json
var defaultConfig embed.FS

func init() {
	configDir, err := getConfigDir()
	if err != nil {
		log.Fatalf("failed to get config dir: %v", err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("failed to create config dir: %v", err)
	}

	dataDir, err := getDataDir()
	if err != nil {
		log.Fatalf("failed to get data dir: %v", err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	configFile := path.Join(configDir, "config.json")
	var configData []byte
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		configData, err = os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
	} else {
		configData, err = defaultConfig.ReadFile("default-config.json")
		if err != nil {
			log.Fatalf("failed to read default config file: %v", err)
		}
		if err := os.WriteFile(configFile, configData, 0644); err != nil {
			log.Fatalf("failed to write config file: %v", err)
		}
	}

	if err := json.Unmarshal(configData, &Config); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	Config.ConfigDir = configDir
	Config.DataDir = dataDir
}
