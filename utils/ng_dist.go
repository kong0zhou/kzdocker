package utils

import (
	"kzdocker/log"
	"net/http"
	"os"
	"strings"
)

// DistHandle angular的静态文件服务器
type DistHandle struct {
	path string
}

// NewDistHandle 工厂函数
func NewDistHandle(p string) *DistHandle {
	return &DistHandle{p}
}

func (t *DistHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	filePath := t.path + upath
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		filePath = t.path + `/index.html`
		http.ServeFile(w, r, filePath)
	} else if err == nil {
		http.ServeFile(w, r, filePath)
	} else {
		log.Error(err.Error())
		return
	}
}
