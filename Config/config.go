package Config

import (
	"net/url"
	"time"
)

var Version string
var BestdoriAPIUrl *url.URL
var BestdoriFanMadeVersion int
var LastUpdate time.Time

func InitConfig() {
	Version = "2.0.0"
	BestdoriAPIUrl, _ = url.Parse("http://202.182.126.173:21104/")
	BestdoriFanMadeVersion = 1
}

func SetLastUpdate(time time.Time) {
	LastUpdate = time
}
