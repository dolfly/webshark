package actions

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/dolfly/webshark/pkg/sharkd"
	"github.com/dolfly/webshark/pkg/sharkd/shark"
	"github.com/gin-gonic/gin"
)

const CAPTURE_PATH = "/captures/"

var sharkclients sync.Map

func getLoaded() []string {
	items := []string{}
	sharkclients.Range(func(key any, val any) bool {
		items = append(items, key.(string))
		return true
	})
	return items
}

func getSharkdClient(capture string) (sharkcli *shark.SharkdClient) {
	key := strings.ReplaceAll(capture, CAPTURE_PATH, "")
	key = strings.TrimPrefix(key, "/")
	val, ok := sharkclients.LoadAndDelete(key)
	if ok {
		sharkcli = val.(*shark.SharkdClient)
	} else {
		sharkcli = shark.NewSharkdClient(sharkd.SockPath())
	}
	sharkclients.Store(key, sharkcli)
	if capture != "" {
		sharkcli.Load(shark.LoadParam{File: "/tmp/test.pcapng"})
	}
	return
}

func ActionUpload(c *gin.Context) {
}

func intval(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func boolval(s string) bool {
	if s == "yes" {
		return true
	} else {
		return false
	}
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ActionJson(c *gin.Context) {
	sharkcli := getSharkdClient(c.DefaultQuery("capture", ""))
	var mmap = map[string]func(){
		"complete": func() {
			r, err := sharkcli.Complete(&shark.CompleteParam{
				Field: c.DefaultQuery("field", ""),
				Pref:  c.DefaultQuery("pref", ""),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"tap": func() {
			r, err := sharkcli.Tap(&shark.TapParam{
				Tap0: c.DefaultQuery("tap", ""),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"frame": func() {
			r, err := sharkcli.Frame(&shark.FrameParam{
				Frame: intval(c.DefaultQuery("frame", "0")),
				Bytes: boolval(c.DefaultQuery("bytes", "")),
				Proto: boolval(c.DefaultQuery("proto", "")),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"frames": func() {
			r, err := sharkcli.Frames(&shark.FramesParam{
				Filter: c.DefaultQuery("filter", ""),
				Skip:   intval(c.DefaultQuery("skip", "")),
				Limit:  intval(c.DefaultQuery("limit", "")),
				Refs:   c.DefaultQuery("refs", ""),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"check": func() {
			r, err := sharkcli.Check(&shark.CheckParam{
				Field:  c.DefaultQuery("field", ""),
				Filter: c.DefaultQuery("filter", ""),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"analyse": func() {
			r, err := sharkcli.Analyse()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"intervals": func() {
			r, err := sharkcli.Intervals(&shark.IntervalsParam{
				Interval: intval(c.DefaultQuery("interval", "0")),
				Filter:   c.DefaultQuery("filter", ""),
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"status": func() {
			r, err := sharkcli.Status()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"info": func() {
			r, err := sharkcli.Info()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"err":    1,
					"errstr": err.Error(),
				})
				return
			}
			data, _ := json.Marshal(r)
			c.Data(http.StatusOK, "application/json", data)
		},
		"files": func() {
			loaded := getLoaded()
			c.JSON(http.StatusOK, gin.H{
				"files": []gin.H{
					{
						"name": filepath.Join(CAPTURE_PATH, "test.pcapng"),
						"size": 100,
						"status": gin.H{
							"online": contains(loaded, "test.pcapng"),
						},
					},
				},
				"pwd": ".",
			})
		},
	}
	method := c.DefaultQuery("req", "")
	mfunc, ok := mmap[method]
	if !ok {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	mfunc()
}
