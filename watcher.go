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
)

// Watcher 热编译功能
type Watcher struct {
	// 项目的源码目录
	Dir string

	// 不同类型日志的输出通道
	Info *log.Logger
	Erro *log.Logger
	Succ *log.Logger

	// 运行 HTTP 服务的函数体
	Serve func() error

	// 执行编译的函数体
	Build func() error

	builded time.Time
}

// Watch 热编译
func Watch(src, base string, info, erro, succ *log.Logger) error {
	u, err := url.Parse(base)
	if err != nil {
		return err
	}

	// 预览文件输出至一个临时目录
	dest := filepath.Join(os.TempDir(), "blogit")

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
	path := u.Path

	w := &Watcher{
		Dir: src,

		Info: info,
		Erro: erro,
		Succ: succ,

		Serve: func() error {
			return Serve(dest, addr, path, info)
		},

		Build: func() error {
			return Build(src, dest, base)
		},

		builded: time.Now(),
	}

	if err := w.Build(); err != nil {
		return err
	}

	return w.Watch()
}

// Watch 监视变化并进行编译
func (w *Watcher) Watch() error {
	go func() {
		if err := w.Serve(); err != nil {
			w.erro(err)
			return
		}
	}()

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

			if time.Now().Sub(w.builded) <= 1*time.Second {
				w.info("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			w.info("触发事件：", event, "，开始重新编译！")

			go func() {
				if err = w.Build(); err != nil {
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
		if info.IsDir() || ext == ".md" || ext == ".yaml" || ext == ".yml" {
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
