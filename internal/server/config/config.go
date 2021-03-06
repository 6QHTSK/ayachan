package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

var Config *YamlConfig
var Version string
var BestdoriAPIUrl *url.URL

type YamlConfig struct {
	RunAddr  string             `yaml:"run-addr"`
	Debug    bool               `yaml:"debug"`
	API      YamlConfigAPI      `yaml:"api"`
	Database YamlConfigDatabase `yaml:"database"`
}

type YamlConfigAPI struct {
	BestdoriAPI string `yaml:"bestdori-api"`
}

type YamlConfigDatabase struct {
	Mysql          string `yaml:"mysql"`
	MeiliSearch    string `yaml:"meilisearch"`
	MeiliSearchKey string `yaml:"meilisearch-key"`
}

func NewYamlConfig() *YamlConfig {
	return &YamlConfig{
		RunAddr: "0.0.0.0:8080",
		Debug:   true,
		API: YamlConfigAPI{
			BestdoriAPI: "http://127.0.0.1:21104",
		},
		Database: YamlConfigDatabase{
			Mysql:          "user:password@/dbname",
			MeiliSearch:    "http://127.0.0.1:7700",
			MeiliSearchKey: "",
		},
	}
}

func init() {
	Version = "2.1.0"
	Config = NewYamlConfig()
	var err error
	if envEnabled() {
		initEnv()
	} else {
		yamlFile, err := ioutil.ReadFile("conf.yaml")
		if err != nil {
			yamlConfig, _ := yaml.Marshal(Config)
			err = ioutil.WriteFile("conf.yaml", yamlConfig, 0666)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("conf.yaml not found, one is generate!")
			os.Exit(0)
		}
		err = yaml.Unmarshal(yamlFile, Config)
		if err != nil {
			log.Fatal("Check the conf.yaml, Cannot Read!")
		}
	}
	BestdoriAPIUrl, err = url.Parse(Config.API.BestdoriAPI)
	if err != nil {
		log.Fatal("Cannot parse BestdoriAPI!")
	}
}
