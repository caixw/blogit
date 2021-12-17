// SPDX-License-Identifier: MIT

package filesystem

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/caixw/blogit/v2/internal/vars"
)

// FileServer 以 fsys 为根目录作为静态文件服务
//
// erro 在出错时日志的输出通道，可以为空，表示输出到 log.Default()；
func FileServer(fsys fs.FS, erro *log.Logger) http.Handler {
	if erro == nil {
		erro = log.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p != "" && p[0] == '/' {
			p = p[1:]
		}

		if p == "" || p[len(p)-1] == '/' {
			p += vars.IndexFilename
		}

	STAT:
		stat, err := fs.Stat(fsys, p)
		if err != nil {
			printError(erro, err, w)
			return
		}
		if stat.IsDir() {
			p = path.Join(p, vars.IndexFilename)
			goto STAT
		}

		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			printError(erro, err, w)
			return
		}

		buf := bytes.NewReader(data)
		http.ServeContent(w, r, filepath.Base(p), stat.ModTime(), buf)
	})
}

func printError(erro *log.Logger, err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, fs.ErrPermission):
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	case errors.Is(err, fs.ErrNotExist):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		erro.Println(err)
	}
}
