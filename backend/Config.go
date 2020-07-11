package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Log                 string `yaml:"log"`
	LogLevel            string `yaml:"log_level"`
	Proxy               string `yaml:"proxy"`
	UserAgent           string `yaml:"user_agent"`
	UseImageCache       bool   `yaml:"use_image_cache"`
	MaxCacheDbSize      int    `yaml:"max_cache_db_size"`
	CacheDbDir          string `yaml:"cache_db_dir"`
	CacheImageDir       string `yaml:"cache_image_dir"`
	ListenAddress       string `yaml:"listen_address"`
	HttpPort            int    `yaml:"http_port"`
	HttpsSupport        bool   `yaml:"https_support"`
	RedirectToHttps     bool   `yaml:"redirect_to_https"`
	RedirectToHttpsCode int    `yaml:"redirect_to_https_code"`
	HttpsPort           int    `yaml:"https_port"`
	CertFile            string `yaml:"cert_file"`
	KeyFile             string `yaml:"key_file"`
}

func LoadConfig(file string) *Configuration {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	configuration := Configuration{
		Log:                 "stdout",
		LogLevel:            "info",
		Proxy:               "",
		UserAgent:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		UseImageCache:       true,
		MaxCacheDbSize:      1024 * 1024 * 100, // 100MB
		CacheDbDir:          "cache/db",
		CacheImageDir:       "cache/image",
		ListenAddress:       "0.0.0.0",
		HttpPort:            80,
		HttpsSupport:        false,
		RedirectToHttps:     true,
		RedirectToHttpsCode: 302,
		HttpsPort:           443,
		CertFile:            "./cert.cert",
		KeyFile:             "./cert.key",
	}
	err = yaml.Unmarshal(data, &configuration)

	out, err := yaml.Marshal(configuration)
	fmt.Printf("%s\n", string(out))

	if err != nil {
		panic(err)
	}
	return &configuration
}

func (_this *Configuration) GetLogFile() io.WriteCloser {
	if _this.Log == "stdout" {
		return os.Stdout
	} else {
		logFile, err := os.Open(_this.Log)
		if err != nil {
			panic(err)
		}
		return logFile
	}
}

func (_this *Configuration) PrepareDirs() {
	_ = os.MkdirAll(_this.CacheDbDir, 0777)
	_ = os.MkdirAll(_this.CacheImageDir, 0777)
}
