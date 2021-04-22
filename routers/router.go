package routers

import (
	"logging-helper/api"
	"logging-helper/middleware/log"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化gin入口，路由信息
func SetupRouter() *gin.Engine {
	router := gin.New()
	if err := log.InitLogger(); err != nil {
		panic(err)
	}
	// router.Use(log.GinLogger(log.Logger),
	// 	log.GinRecovery(log.Logger, true))

	router.GET("/health", api.Health)

	v1 := router.Group("/log/api/v1")
	{
		v1.GET("/all/task/:task", api.QueryAllLog)
		v1.GET("/sse/task/:task", api.QuerySseLog)
	}

	return router
}
