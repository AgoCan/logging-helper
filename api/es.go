package api

import (
	"io"
	"logging-helper/service"
	"logging-helper/utils/sse"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryAllLog(c *gin.Context) {
	task := c.Params.ByName("task")
	if task == "" {
		c.JSON(http.StatusOK, gin.H{"task": task, "status": "no value"})
	}
	service := service.Elastic{}
	res := service.Query(task)
	c.JSON(http.StatusOK, gin.H{"task": task, "status": res})
}

func QuerySseLog(c *gin.Context) {
	task := c.Params.ByName("task")
	if task == "" {
		c.JSON(http.StatusOK, gin.H{"task": task, "status": "no value"})
		return
	}
	stream := sse.NewServer()
	service := service.Elastic{}
	clientChan := make(sse.ClientChan)

	// Send new connection to event server
	stream.NewClients <- clientChan

	defer func() {
		// Send closed connection to event server
		stream.ClosedClients <- clientChan
	}()
	go service.QueryBySse(task, stream)

	c.Stream(func(w io.Writer) bool {

	queryLoop:
		for {
			select {
			case msg := <-stream.Message:
				c.SSEvent("message", msg)
				return true
			case <-c.Request.Context().Done():
				stream.Stop = 1
				break queryLoop
			}
		}
		return false
	})
}
