package web

import (
	"embed"
	_ "embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/static
var static embed.FS

//go:embed templates
var tmpls embed.FS

//go:embed dist/index.html
var index []byte

//go:embed dist/config.html
var config []byte

//go:embed dist/ipmap.html
var ipmap []byte

//go:embed dist/upload.html
var upload []byte

//go: embed dist/favicon.ico
var favicon []byte

type embedFS struct {
	embed.FS
	path string
}
type embedFile struct {
	io.Seeker
	fs.File
}

func (*embedFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}
func (fs *embedFS) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	name = strings.Split(name, "?")[0]
	fullName := filepath.Join(fs.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := fs.FS.Open(fullName)
	ef := &embedFile{
		File: file,
	}
	return ef, err
}
func SetTemplate(r *gin.Engine) error {
	t, err := template.ParseFS(tmpls, "templates/*.tmpl", "templates/*/*.tmpl")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(t)
	return nil
}

func SetRoute(r *gin.Engine, prefix string) {
	r.StaticFS("/static", Static())
	pgroup := r.Group(prefix)
	{
		pgroup.GET("/", Index())
		pgroup.GET("/config", Config())
		pgroup.GET("/ipmap", Ipmap())
		pgroup.GET("/upload", Upload())
	}
}

// Static Static
func Static() http.FileSystem {
	if err := recover(); err != nil {
		//fmt.Println(err)
	}
	return &embedFS{
		static,
		"dist/static",
	}
}

// Index Index
func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", index)
	}
}

func Config() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", config)
	}
}

func Ipmap() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", ipmap)
	}
}

func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", upload)
	}
}

// Favicon Favicon
func Favicon() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "image/ico", favicon)
	}
}
