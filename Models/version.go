package Models

import "ayachan/Config"

type APIVersion struct {
	Version string `json:"version"`
}

func (version *APIVersion) GetVersion() {
	version.Version = Config.Version
}
