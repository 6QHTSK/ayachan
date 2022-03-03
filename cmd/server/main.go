package main

import (
	"flag"
	"fmt"
	"github.com/6QHTSK/ayachan/internal/server"
	"github.com/6QHTSK/ayachan/internal/server/config"
)

var showVer bool
var runAddr string

func init() {
	flag.BoolVar(&showVer, "v", false, "查看版本号")
	flag.StringVar(&runAddr, "a", config.Config.RunAddr, "运行地址")
}

func main() {
	flag.Parse()
	if showVer {
		fmt.Println(config.Version)
	} else {
		router := server.InitRouter()

		server.InitAPI(router)
		_ = router.Run(runAddr)
	}
}
