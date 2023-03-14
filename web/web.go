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

//go:embed dist/*.html dist/*.ico
var webfile embed.FS

func SetTemplate(r *gin.Engine) error {
	t, err := template.ParseFS(tmpls, "templates/*.tmpl", "templates/*/*.tmpl")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(t)
	return nil
}

func SetRoute(r *gin.Engine, prefix string) {
	r.StaticFS("/static", &FS{static, "dist/static"})
	r.StaticFS("/webshark", &FS{webfile, "dist"})
}

type FS struct {
	embed.FS
	path string
}
type File struct {
	io.Seeker
	fs.File
}

func (*File) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}
func (fs *FS) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	name = strings.Split(name, "?")[0]
	fullName := filepath.Join(fs.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := fs.FS.Open(fullName)
	ef := &File{
		File: file,
	}
	return ef, err
}
