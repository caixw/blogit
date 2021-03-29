// SPDX-License-Identifier: MIT

package preview

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/caixw/blogit"
	"github.com/caixw/blogit/filesystem"
)

// options 启动服务的参数选项
type options struct {
	// 项目的源码目录
	// 如果为空，则会采用 ./ 作为默认值。
	source string
	srcFS  fs.FS

	// 项目编译后的输出地址
	// 如果为空，则会要用 filesystem.Memory() 作为默认值。
	dest string

	// 如果指定了此值，那么表示要替换 conf.yaml 中的 url
	url  string
	addr string
	path string

	// HTTPS 模式下的证书
	cert string
	key  string

	info *log.Logger
	erro *log.Logger
	succ *log.Logger

	b       *blogit.Builder
	srv     *http.Server
	builded time.Time

	stop chan struct{}
}

func (o *options) sanitize() error {
	if o.info == nil {
		o.info = log.New(os.Stdout, "", log.LstdFlags)
	}

	if o.succ == nil {
		o.succ = log.New(os.Stdout, "", log.LstdFlags)
	}

	if o.erro == nil {
		o.erro = log.New(os.Stderr, "", log.LstdFlags)
	}

	if o.source == "" {
		o.source = "./"
	}
	o.srcFS = os.DirFS(o.source)

	u, err := url.Parse(o.url)
	if err != nil {
		return err
	}

	o.addr = u.Port()
	if o.addr == "" {
		switch scheme := strings.ToLower(u.Scheme); scheme {
		case "https":
			o.addr = ":443"
		case "http", "":
			o.addr = ":80"
		default:
			return fmt.Errorf("不支持协议：%s", scheme)
		}
	} else {
		o.addr = ":" + o.addr
	}

	o.path = u.Path
	if o.path == "" || o.path[0] != '/' {
		o.path = "/" + o.path
	}

	var dest filesystem.WritableFS
	if o.dest == "" {
		dest = filesystem.Memory()
	} else {
		dest = filesystem.Dir(o.dest)
	}
	o.b = blogit.NewBuilder(dest, o.erro)

	o.srv = &http.Server{Addr: o.addr, Handler: o.initServer()}

	o.stop = make(chan struct{}, 1)

	return nil
}

func (o *options) initServer() http.Handler {
	var h http.Handler = o.b

	if o.info != nil {
		o.info.Println("启动服务：", o.addr)

		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o.info.Printf("访问 %s\n", r.URL.String())
			o.b.ServeHTTP(w, r)
		})
	}
	return http.StripPrefix(o.path, h)
}

func (o *options) build() error {
	return o.b.Rebuild(o.srcFS, o.url)
}

func (o *options) watch() error {
	if err := o.sanitize(); err != nil {
		return err
	}

	go func() {
		if err := o.serve(); !errors.Is(err, http.ErrServerClosed) {
			o.erro.Println(err)
		}
		o.builded = time.Now()
	}()

	if err := o.build(); err != nil {
		o.erro.Println(err)
	}

	watcher, err := o.getWatcher()
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod { // 忽略 CHMOD 事件
				continue
			}

			if time.Since(o.builded) <= time.Second {
				o.info.Println("watcher.Events:更新太频繁，该监控事件被忽略:", event)
				continue
			}

			o.info.Println("触发事件：", event, "，开始重新编译！")

			go func() {
				if err = o.build(); err != nil {
					o.erro.Println(err)
					return
				}
				o.builded = time.Now()
				o.succ.Println("重新编译成功")
			}()
		case err := <-watcher.Errors:
			o.erro.Println(err)
			return err
		case <-o.stop:
			return http.ErrServerClosed
		}
	}
}

func (o *options) close() error {
	o.stop <- struct{}{}
	return o.srv.Close()
}

func (o *options) serve() error {
	if o.cert != "" && o.key != "" {
		return o.srv.ListenAndServeTLS(o.cert, o.key)
	}
	return o.srv.ListenAndServe()
}

func (o *options) getWatcher() (*fsnotify.Watcher, error) {
	paths := make([]string, 0, 10)
	err := filepath.Walk(o.source, func(p string, d os.FileInfo, err error) error {
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
