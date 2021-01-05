// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
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

// Version 版本号
const Version = "0.1.0"

// Build 编译并输出内容
func Build(src, target string) error {
	b, err := NewBuilder(src)
	if err != nil {
		return err
	}
	return b.Dump(target)
}

// Serve 运行服务
func Serve(src, addr, path string) error {
	b, err := NewBuilder(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
func ServeTLS(src, addr, path, cert, key string) error {
	b, err := NewBuilder(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServeTLS(addr, cert, key, nil)
}

func Watch(src, addr, path string) error {
	b, err := NewBuilder(src)
	if err != nil {
		return err
	}

	watcher, err := getWatcher(src)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Println("watcher.Events:忽略 CHMOD 事件:", event)
				continue
			}

			if time.Now().Sub(b.Builded) <= 1*time.Second { // 已经记录
				log.Println("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			log.Println("watcher.Events:触发加载事件:", event)

			go func() {
				if err = watcher.Close(); err != nil {
					log.Println(err)
					return
				}
				if err = Watch(src, addr, path); err != nil {
					log.Println(err)
					return
				}
			}()
		case err := <-watcher.Errors:
			log.Println(err)
			return err
		}
	}
}

// NewBuilder 将数据转换成 Builder 对象
func NewBuilder(dir string) (*builder.Builder, error) {
	d, err := data.Load(dir)
	if err != nil {
		return nil, err
	}

	return builder.Build(d)
}

func getWatcher(dir string) (*fsnotify.Watcher, error) {
	paths := make([]string, 0, 10)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if err == nil && (info.IsDir() || ext == ".md" || ext == ".yaml" || ext == ".yml") {
			paths = append(paths, path)
		}
		return err
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
