package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

var Config *YamlConfig
var BestdoriAPIUrl *url.URL

type YamlConfig struct {
	Mysql          string `yaml:"mysql"`
	MeiliSearch    string `yaml:"meilisearch"`
	MeiliSearchKey string `yaml:"meilisearch-key"`
	BestdoriAPI    string `yaml:"bestdori-api"`
}

func NewYamlConfig() *YamlConfig {
	return &YamlConfig{
		BestdoriAPI:    "http://127.0.0.1:21104",
		Mysql:          "user:password@/dbname",
		MeiliSearch:    "http://127.0.0.1:7700",
		MeiliSearchKey: "",
	}
}

func init() {
	Config = NewYamlConfig()
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
	BestdoriAPIUrl, err = url.Parse(Config.BestdoriAPI)
	if err != nil {
		log.Fatal("Cannot parse BestdoriAPI!")
	}
}
