package main

import (
	"ayachanV2/Config"
	"ayachanV2/Router"
	"fmt"
)

// @title ayachan API
// @version 2.0
// @description api 计算Bestdori谱面难度，获得Bestdori数据，常见Bandori谱面格式转换等

// @contact.name 6QHTSK

// @license.name Apache 2.0
// @license.url https://mit-license.org/

// @host 127.0.0.1:8080
// @BasePath /v2

func main() {
	fmt.Println("Hello World!")
	Config.InitConfig()

	router := Router.InitRouter()

	Router.InitSwaggerDoc(router)
	Router.InitAPIV2(router)

	_ = router.Run("0.0.0.0:8080")
}
