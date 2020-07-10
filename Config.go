package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Log            string `yaml:"log"`
	LogLevel       string `yaml:"log_level"`
	Proxy          string `yaml:"proxy"`
	UserAgent      string `yaml:"user_agent"`
	MaxCacheDbSize int    `yaml:"max_cache_db_size"`
	CacheDbDir     string `yaml:"cache_db_dir"`
	CacheImageDir  string `yaml:"cache_image_dir"`
}

func LoadConfig(file string) *Configuration {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	configuration := Configuration{
		Log:            "stdout",
		LogLevel:       "info",
		Proxy:          "",
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		MaxCacheDbSize: 1024 * 1024 * 100, // 100MB
		CacheDbDir:     "cache/db",
		CacheImageDir:  "cache/image",
	}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		panic(err)
	}
	return &configuration
}
