package main

import (
	"fmt"
	_ "logging-helper/config"
	"logging-helper/routers"
)

var (
	err error
)

func main() {

	// 调用路由组
	router := routers.SetupRouter()

	err = router.Run(":9000")
	if err != nil {
		fmt.Println(err)
	}
}
