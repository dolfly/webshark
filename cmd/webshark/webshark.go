package main

import (
	"github.com/dolfly/webshark/api/actions"
	"github.com/dolfly/webshark/web"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	web.SetRoute(r, "/webshark")
	web.SetTemplate(r)
	webshark := r.Group("/webshark")
	{
		webshark.POST("/upload", actions.ActionUpload)
		webshark.GET("/json", actions.ActionJson)
	}
	r.Run(":21405")
}
