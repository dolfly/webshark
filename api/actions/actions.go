package actions

import (
	"net/http"

	"github.com/dolfly/webshark/pkg/sharkd"
	"github.com/gin-gonic/gin"
)

var (
	sharkdcli = sharkd.NewSharkdClient()
)

func ActionUpload(c *gin.Context) {
}

func ActionJson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
