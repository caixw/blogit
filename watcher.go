// SPDX-License-Identifier: MIT

package blogit

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/vars"
)

// Watcher 热编译功能
type Watcher struct {
	// 项目的源码目录
	Dir string

	BaseURL string
	addr    string
	path    string

	// 不同类型日志的输出通道
	Info *log.Logger
	Erro *log.Logger
	Succ *log.Logger

	builder *builder.Builder
	builded time.Time
}

// Watch 热编译
//
// src 源码目录，该目录下的内容一量修改，就会重新编译；
// base 网站的根地址，会替换配置文件中的 URL，一般为 http://localhost，同时也会作为服务的监听地址。
func Watch(src, base string, info, erro, succ *log.Logger) error {
	u, err := url.Parse(base)
	if err != nil {
		return err
	}

	addr := u.Port()
	if addr == "" {
		if scheme := strings.ToLower(u.Scheme); scheme == "https" {
			addr = ":443"
		} else if scheme == "http" {
			addr = ":http"
		} else {
			return fmt.Errorf("不支持协议：%s", scheme)
		}
	} else {
		addr = ":" + addr
	}

	w := &Watcher{
		Dir: src,

		BaseURL: base,
		addr:    addr,
		path:    u.Path,

		Info: info,
		Erro: erro,
		Succ: succ,

		builder: &builder.Builder{},

		builded: time.Now(),
	}

	return w.Watch()
}

// Watch 监视变化并进行编译
func (w *Watcher) Watch() error {
	go func() {
		if err := serve(w.builder, w.Dir, w.addr, w.path, "", "", w.Info); err != nil {
			w.erro(err)
			return
		}
	}()

	if err := w.builder.Build(w.Dir, w.BaseURL); err != nil {
		w.erro(err)
	}

	watcher, err := w.getWatcher()
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod || // 忽略 CHMOD 事件
				strings.ToLower(filepath.Ext(event.Name)) == ".xml" { // 忽略对 xml 文件的写操作
				continue
			}

			if time.Since(w.builded) <= time.Second {
				w.info("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			w.info("触发事件：", event, "，开始重新编译！")

			go func() {
				if err = w.builder.Build(w.Dir, w.BaseURL); err != nil {
					w.erro(err)
					return
				}
				w.builded = time.Now()
				w.succ("重新编译成功")
			}()
		case err := <-watcher.Errors:
			w.info(err)
			return err
		}
	}
}

func (w *Watcher) getWatcher() (*fsnotify.Watcher, error) {
	paths := make([]string, 0, 10)
	err := filepath.Walk(w.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()

		if name != "." && name[0] == '.' { // 忽略隐藏文件
			return filepath.SkipDir
		}

		ext := strings.ToLower(filepath.Ext(name))
		if info.IsDir() || ext == vars.MarkdownExt || ext == ".yaml" || ext == ".yml" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if err := watcher.Add(p); err != nil {
			return nil, err
		}
	}

	return watcher, nil
}

func (w *Watcher) info(v ...interface{}) {
	if w.Info != nil {
		w.Info.Println(v...)
	}
}

func (w *Watcher) erro(v ...interface{}) {
	if w.Erro != nil {
		w.Erro.Println(v...)
	}
}

func (w *Watcher) succ(v ...interface{}) {
	if w.Succ != nil {
		w.Succ.Println(v...)
	}
}
