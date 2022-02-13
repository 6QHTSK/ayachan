package Config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/url"
	"time"
)

var Config *YamlConfig
var Version string
var BestdoriAPIUrl *url.URL
var BestdoriFanMadeVersion int
var LastUpdate time.Time

type YamlConfig struct {
	RunAddr  string             `yaml:"run-addr"`
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
	Version = "2.0.1"
	Config = NewYamlConfig()
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		yamlConfig, err := yaml.Marshal(Config)
		err = ioutil.WriteFile("conf.yaml", yamlConfig, 0666)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Fatal("conf.yaml not found, one is generate!")
	}
	err = yaml.Unmarshal(yamlFile, Config)
	if err != nil {
		log.Fatal("Check the conf.yaml, Cannot Read!")
	}
	BestdoriAPIUrl, err = url.Parse(Config.API.BestdoriAPI)
	if err != nil {
		log.Fatal("Cannot parse BestdoriAPI!")
	}
	BestdoriFanMadeVersion = 1
}

func SetLastUpdate(time time.Time) {
	LastUpdate = time
}
