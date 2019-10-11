package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

const (
	kConfigFileName = "./config.json"
)

type LogInfo struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

type (
	WebInfo struct {
		WebName       string
		Homepage      string
		SearchApi     string
		Method        string
		ContentReg    string
		ListReg       string
		SearchReg     string
		SearchPageReg string
	}
)

type Config struct {
	Log          LogInfo   `json:"log"`
	Listen       string    `json:"listen"`
	NovelWebInfo []WebInfo `json:"novelWebInfo"`
	closeChan    chan struct{}
	doneChan     chan struct{}
}

var config *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config, err = newConfig()
	})

	return config, err
}

func newConfig() (*Config, error) {
	cfg := &Config{}
	cfg.closeChan = make(chan struct{})
	cfg.doneChan = make(chan struct{})
	err := cfg.load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Close() {
	close(c.closeChan)
	<-c.doneChan
}

func (c *Config) load() error {
	content, err := ioutil.ReadFile(kConfigFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		return err
	}

	return nil
}
