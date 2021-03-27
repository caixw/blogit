// SPDX-License-Identifier: MIT

package blogit

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/utils"
	"github.com/caixw/blogit/internal/vars"
)

// Watcher 热编译功能
type Watcher struct {
	dir     string // 项目的源码目录
	baseURL string
	addr    string
	path    string
	cert    string
	key     string

	// 不同类型日志的输出通道
	infoLog *log.Logger
	erroLog *log.Logger
	succLog *log.Logger

	builder *builder.Builder
	builded time.Time
}

// Watch 热编译
//
// src 源码目录，该目录下的内容一量修改，就会重新编译；
// base 网站的根地址，会替换配置文件中的 URL，
// 一般为 http://localhost，同时也会作为服务的监听地址。
// 如果是以 https:// 开头的，那么需要提供 cert 和 key 两个参数；
// cert 和 key 表示 https 模式下对应的证书；
// info, erro, succ 为各类型的日志输出通道，可以为空，表示该类型的信息将会被忽略。
func Watch(src, base, cert, key string, info, erro, succ *log.Logger) (*Watcher, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	addr := u.Port()
	scheme := strings.ToLower(u.Scheme)
	if addr == "" {
		if scheme == "https" {
			addr = ":443"
		} else if scheme == "http" {
			addr = ":http"
		} else {
			return nil, fmt.Errorf("不支持协议：%s", scheme)
		}
	} else {
		addr = ":" + addr
	}

	if scheme == "https" && (!utils.FileExists(cert) || !utils.FileExists(key)) {
		return nil, errors.New("HTTPS 模式但是证书不存在")
	}

	return &Watcher{
		dir:     src,
		baseURL: base,
		addr:    addr,
		path:    u.Path,
		cert:    cert,
		key:     key,

		infoLog: info,
		erroLog: erro,
		succLog: succ,

		builder: &builder.Builder{},
		builded: time.Now(),
	}, nil
}

// Watch 监视变化并进行编译
func (w *Watcher) Watch() error {
	go func() {
		if err := serve(w.builder, w.dir, w.addr, w.path, w.cert, w.key, w.infoLog); err != nil {
			w.erro(err)
			return
		}
	}()

	if err := w.builder.Build(w.dir, w.baseURL); err != nil {
		w.erro(err)
	}

	watcher, err := w.getWatcher()
	if err != nil {
		return err
	}

	w.info("启动服务：", w.baseURL)
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
				if err = w.builder.Build(w.dir, w.baseURL); err != nil {
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
	err := filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
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
	if w.infoLog != nil {
		w.infoLog.Println(v...)
	}
}

func (w *Watcher) erro(v ...interface{}) {
	if w.erroLog != nil {
		w.erroLog.Println(v...)
	}
}

func (w *Watcher) succ(v ...interface{}) {
	if w.succLog != nil {
		w.succLog.Println(v...)
	}
}
