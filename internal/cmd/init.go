// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"io/fs"
	"path"
	"time"

	"github.com/issue9/cmdopt"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

var initFS *flag.FlagSet

// initInit 注册 init 子命令
func initInit(opt *cmdopt.CmdOpt) {
	initFS = opt.New("init", "初始化新的博客内容\n", initF)
}

func initF(w io.Writer) error {
	if initFS.NArg() != 1 {
		erro.println("必须指定目录")
		return nil
	}

	wfs, err := getWD()
	if err != nil {
		erro.println(err)
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
		erro.println(err)
		return nil
	}
	succ.println("创建了文件:", vars.ConfYAML)

	// tags.yaml
	tags := []*loader.Tag{
		{
			Slug:    "default",
			Title:   "默认",
			Content: "这是默认的标签",
		},
	}
	if err := writeYAML(wfs, vars.TagsYAML, tags); err != nil {
		erro.println(err)
		return nil
	}
	succ.println("创建了文件:", vars.TagsYAML)

	// themes
	theme := &loader.Theme{
		URL:         "https://example.com",
		Description: "description",
	}
	p := path.Join(vars.ThemesDir, "default", "theme.yaml")
	if err := writeYAML(wfs, p, theme); err != nil {
		erro.println(err)
		return nil
	}
	succ.println("创建了主题文件:", p)

	p = path.Join(vars.PostsDir, time.Now().Format("2006"), "post1.md")
	if err := wfs.WriteFile(p, []byte(content), fs.ModePerm); err != nil {
		erro.println(err)
		return nil
	}
	succ.println("创建了文章:", p)

	return nil
}

func writeYAML(wfs filesystem.WritableFS, path string, v interface{}) error {
	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return wfs.WriteFile(path, bs, fs.ModePerm)
}
