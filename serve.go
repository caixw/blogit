// SPDX-License-Identifier: MIT

package blogit

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Server 表示一个长期运行的服务对象的基本操作
type Server interface {
	// 执行服务
	//
	// 这是一个阻塞函数，直到 Close 关闭服务才会返回 http.ErrServerClosed。
	Serve() error

	// 关闭服务
	Close() error
}

type server struct {
	*Options
}

// watcher 热编译功能
type watcher struct {
	*Options
	dir     string
	builded time.Time
	stop    chan struct{}
}

// Serve HTTP 服务对象
func Serve(o *Options) (Server, error) {
	if err := o.init(); err != nil {
		return nil, err
	}
	return &server{Options: o}, nil
}

func (s *server) Serve() error {
	if err := s.build(); err != nil {
		return err
	}
	return s.serve()
}

func (s *server) Close() error {
	return s.srv.Close()
}

// Watch 热编译服务对象
//
// dir 是指需要监视的目录，当该目录下的文件发生改变时，即会重新编译项目；
func Watch(dir string, o *Options) (Server, error) {
	if err := o.init(); err != nil {
		return nil, err
	}

	return &watcher{
		Options: o,
		dir:     dir,
		builded: time.Now(),
		stop:    make(chan struct{}, 1),
	}, nil
}

func (w *watcher) Serve() error {
	go func() {
		if err := w.serve(); !errors.Is(err, http.ErrServerClosed) {
			w.erro(err)
		}
	}()

	if err := w.build(); err != nil {
		w.erro(err)
	}

	watcher, err := w.getWatcher()
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod { // 忽略 CHMOD 事件
				continue
			}

			if time.Since(w.builded) <= time.Second {
				w.info("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			w.info("触发事件：", event, "，开始重新编译！")

			go func() {
				if err = w.build(); err != nil {
					w.erro(err)
					return
				}
				w.builded = time.Now()
				w.succ("重新编译成功")
			}()
		case err := <-watcher.Errors:
			w.erro(err)
			return err
		case <-w.stop:
			w.srv.Close()
			return http.ErrServerClosed
		}
	}
}

func (w *watcher) Close() error {
	w.stop <- struct{}{}
	return nil
}

func (w *watcher) getWatcher() (*fsnotify.Watcher, error) {
	paths := make([]string, 0, 10)
	err := filepath.Walk(w.dir, func(p string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if name != "." && name[0] == '.' { // 忽略隐藏文件
			return filepath.SkipDir
		}

		paths = append(paths, p)
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
