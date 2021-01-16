// SPDX-License-Identifier: MIT

package blogit

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher 热编译功能
type Watcher struct {
	Dir     string
	Info    *log.Logger
	Erro    *log.Logger
	Succ    *log.Logger
	F       func() error
	builded time.Time
}

// Watch 热编译
func Watch(src, addr, path string, info, erro, succ *log.Logger) error {
	w := &Watcher{
		Dir:  src,
		Info: info,
		Erro: erro,
		Succ: succ,
		F: func() error {
			return Serve(src, addr, path, info)
		},
		builded: time.Now(),
	}

	if err := Build(src); err != nil {
		return err
	}

	return w.Watch()
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

// Watch 监视变化并进行编译
func (w *Watcher) Watch() error {
	go func() {
		if err := w.F(); err != nil {
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
				if err = Build(w.Dir); err != nil {
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
