package api

import (
	"logging-helper/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	service := service.Health{}
	res := service.Status()
	c.JSON(http.StatusOK, res)
}
