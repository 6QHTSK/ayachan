package Models

import "github.com/6QHTSK/ayachan/Config"

type APIVersion struct {
	Version string `json:"version"`
}

func (version *APIVersion) GetVersion() {
	version.Version = Config.Version
}
