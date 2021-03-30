// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"io"
	"io/fs"
	"log"
	"path"
	"time"

	"github.com/issue9/cmdopt"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

var initFS *flag.FlagSet

// InitInit 注册 init 子命令
func InitInit(opt *cmdopt.CmdOpt, succ, erro *log.Logger) {
	initFS = opt.New("init", "初始化新的博客内容\n", initF(succ, erro))
}

func initF(succ, erro *log.Logger) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if initFS.NArg() != 1 {
			erro.Panicln("必须指定目录")
			return nil
		}

		wfs, err := getWD()
		if err != nil {
			erro.Panicln(err)
			return nil
		}

		// conf.yaml
		conf := &loader.Config{
			Title:  "example",
			URL:    "https://example.com",
			Uptime: time.Now(),
			Theme:  "default",
		}
		if err := writeYAML(wfs, vars.ConfYAML, conf); err != nil {
			erro.Panicln(err)
			return nil
		}
		succ.Println("创建了文件:", vars.ConfYAML)

		// tags.yaml
		tags := []*loader.Tag{
			{
				Slug:    "default",
				Title:   "默认",
				Content: "这是默认的标签",
			},
		}
		if err := writeYAML(wfs, vars.TagsYAML, tags); err != nil {
			erro.Panicln(err)
			return nil
		}
		succ.Println("创建了文件:", vars.TagsYAML)

		// themes
		theme := &loader.Theme{
			URL:         "https://example.com",
			Description: "description",
		}
		p := path.Join(vars.ThemesDir, "default", "theme.yaml")
		if err := writeYAML(wfs, p, theme); err != nil {
			erro.Panicln(err)
			return nil
		}
		succ.Println("创建了主题文件:", p)

		p = path.Join(vars.PostsDir, time.Now().Format("2006"), "post1.md")
		if err := wfs.WriteFile(p, []byte(postContent), fs.ModePerm); err != nil {
			erro.Panicln(err)
			return nil
		}
		succ.Println("创建了文章:", p)

		return nil
	}
}

func writeYAML(wfs filesystem.WritableFS, path string, v interface{}) error {
	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return wfs.WriteFile(path, bs, fs.ModePerm)
}
