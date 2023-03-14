package main

import (
	"context"

	"github.com/dolfly/webshark/api/actions"
	"github.com/dolfly/webshark/pkg/sharkd"
	"github.com/dolfly/webshark/web"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sharkd.Start(ctx)
	defer cancel()
	r := gin.Default()
	web.SetRoute(r, "/webshark")
	api := r.Group("/api")
	{
		api.POST("/upload", actions.ActionUpload)
		api.GET("/json", actions.ActionJson)
	}
	r.Run(":21405")
}
