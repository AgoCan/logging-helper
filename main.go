package main

import (
	"logging-helper/routers"
	"logging-helper/service"
)

func main() {
	err := service.InitElasticClinet()
	if err != nil {
		panic(err)
	}

	// 调用路由组
	router := routers.SetupRouter()

	err = router.Run(":9000")
	if err != nil {
		panic(err)
	}

}
