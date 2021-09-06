// SPDX-License-Identifier: MIT

package preview

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/cmd/console"
)

// options 启动服务的参数选项
type options struct {
	p *message.Printer

	// 项目的源码目录
	// 如果为空，则会采用 ./ 作为默认值。
	source string
	srcFS  fs.FS

	// 项目编译后的输出地址
	// 如果为空，则会要用 builder.Memory() 作为默认值。
	dest   string
	destFS blogit.WritableFS

	url  string // 如果指定了此值，那么表示要替换 conf.yaml 中的 url
	addr string
	path string

	// HTTPS 模式下的证书
	cert string
	key  string

	b       *blogit.Builder
	srv     *http.Server
	builded time.Time

	stop chan struct{}
}

func (o *options) sanitize() error {
	o.stop = make(chan struct{}, 1)

	if o.source == "" {
		o.source = "./"
	}
	o.srcFS = os.DirFS(o.source)

	if err := o.parseURL(); err != nil {
		return err
	}

	if o.dest == "" {
		o.destFS = blogit.MemoryFS()
	} else {
		o.destFS = blogit.DirFS(o.dest)
	}

	return nil
}

func (o *options) parseURL() error {
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
			return fmt.Errorf(o.p.Sprintf("preview unknown protocol", scheme))
		}
	} else {
		o.addr = ":" + o.addr
	}

	o.path = u.Path
	if o.path == "" || o.path[0] != '/' {
		o.path = "/" + o.path
	}

	return nil
}

func (o *options) build(erro *console.Logger) (ok bool) {
	if err := o.b.Rebuild(); err != nil {
		if ls, ok := err.(localeutil.LocaleStringer); ok {
			erro.Println(ls.LocaleString(o.p))
		} else {
			erro.Println(err)
		}
		return false
	}

	o.builded = time.Now()
	return true
}

func (o *options) watch(succ, info, erro *console.Logger) error {
	if err := o.sanitize(); err != nil {
		return err
	}

	o.b = &blogit.Builder{
		Src:     o.srcFS,
		Dest:    o.destFS,
		Info:    info.AsLogger(),
		Preview: true,
		BaseURL: o.url,
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info.Println(o.p.Sprintf("visit url", r.URL.String()))
		o.b.Handler(erro.AsLogger()).ServeHTTP(w, r)
	})
	o.srv = &http.Server{Addr: o.addr, Handler: http.StripPrefix(o.path, h)}

	go func() {
		if err := o.serve(info); !errors.Is(err, http.ErrServerClosed) {
			erro.Println(err)
		}
		o.stop <- struct{}{}
	}()

	o.build(erro)

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
				info.Println(o.p.Sprintf("preview ignore", event))
				continue
			}

			info.Println(o.p.Sprintf("preview trigger event", event))

			go func() {
				if !o.build(erro) {
					return
				}
				succ.Println(o.p.Sprintf("preview rebuild success"))
			}()
		case err := <-watcher.Errors:
			erro.Println(err)
			return err
		case <-o.stop:
			return http.ErrServerClosed
		}
	}
}

func (o *options) close() error { return o.srv.Close() }

func (o *options) serve(info *console.Logger) error {
	info.Println(o.p.Sprintf("start server", o.addr))

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

		if isHidden(d.Name()) { // 忽略隐藏文件
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
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

func isHidden(name string) bool {
	return len(name) > 2 && name[0] == '.' && name[1] != '/' && name[1] != os.PathSeparator
}
