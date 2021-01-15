// SPDX-License-Identifier: MIT

package blogit

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/data"
)

// Watcher 热编译功能
type Watcher struct {
	Dir string
	Log *log.Logger
	F   func(http.Handler) error

	b *builder.Builder
}

// Watch 热编译
func Watch(src, addr, path string, l *log.Logger) error {
	w := &Watcher{
		Dir: src,
		Log: l,
		F: func(h http.Handler) error {
			http.Handle(path, h)
			return http.ListenAndServe(addr, nil)
		},
		b: &builder.Builder{},
	}

	if err := w.load(); err != nil {
		return err
	}

	return w.Watch()
}

// Watch 监视变化并进行编译
func (w *Watcher) Watch() error {
	go func() {
		if err := w.F(w.b); err != nil {
			w.Log.Println(err)
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
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				w.Log.Println("watcher.Events:忽略 CHMOD 事件:", event)
				continue
			}

			if time.Now().Sub(w.b.Builded) <= 1*time.Second {
				w.Log.Println("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			w.Log.Println("watcher.Events:触发事件:", event)

			go func() {
				if err := w.load(); err != nil {
					w.Log.Println(err)
					return
				}
			}()
		case err := <-watcher.Errors:
			w.Log.Println(err)
			return err
		}
	}
}

func (w *Watcher) load() error {
	d, err := data.Load(w.Dir)
	if err != nil {
		return err
	}

	return w.b.Load(d)
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
