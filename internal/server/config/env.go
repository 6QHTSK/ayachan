package config

import "os"

func envEnabled() bool {
	return os.Getenv("use_env") == "true"
}

func readEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initEnv() {
	Config = &YamlConfig{
		RunAddr: readEnv("run_addr", Config.RunAddr),
		Debug:   false,
		API: YamlConfigAPI{
			BestdoriAPI: readEnv("bestdori_api", Config.API.BestdoriAPI),
		},
		Database: YamlConfigDatabase{
			Mysql:          readEnv("mysql", Config.Database.Mysql),
			MeiliSearch:    readEnv("meilisearch", Config.Database.MeiliSearch),
			MeiliSearchKey: readEnv("meilisearch_key", Config.Database.MeiliSearchKey),
		},
	}
}
