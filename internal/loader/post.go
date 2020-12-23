// SPDX-License-Identifier: MIT

package loader

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/issue9/sliceutil"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

// 表示 Post.State 的各类值
const (
	StateTop     = "top"   // 置顶
	StateLast    = "last"  // 放在尾部
	StateDraft   = "draft" // 表示为草稿，不会加载此条数据
	StateDefault = ""      // 默认值
)

// 文章是否过时的比较方式
const (
	OutdatedCreated  = "created"
	OutdatedModified = "modified"
)

// Post 表示文章的信息
type Post struct {
	Title     string `yaml:"title"`
	HTMLTitle string `yaml:"-"` // 表示网页中 html>title 中的内容，可能带网站名称的后缀

	Created  time.Time `yaml:"created"`           // 创建时间
	Modified time.Time `yaml:"modified"`          // 修改时间
	Summary  string    `yaml:"summary,omitempty"` // 摘要，同时也作为 meta.description 的内容

	// 关联的标签列表
	//
	// 标签名为各个标签的 slug 值，可以保证其唯一。
	// 最终会被解析到 TagString 中，TagString 会被废弃。
	TagString []string `yaml:"tags"`
	Tags      []*Tag   `yaml:"-"`

	// Outdated 用户记录文章的一个过时情况，可以由以下几种值构成：
	// - created 表示该篇文章以创建时间来计算其是否已经过时，该值也是默认值；
	// - modified 表示该文章以其修改时间来计算其是否已经过时；
	// - none 表示该文章永远不会过时；
	// - 其它任意非空值，表示直接以该字符串当作过时信息展示给用语.
	Outdated string `yaml:"outdated,omitempty"`

	// State 表示文章的状态，有以下四种值：
	// - top 表示文章被置顶；
	// - last 表示文章会被放置在最后；
	// - draft 表示这是一篇草稿，并不会被加地到内存中；
	// - default 表示默认情况，也可以为空，按默认的方式进行处理。
	State string `yaml:"state,omitempty"`

	// 封面地址，可以为空。
	Image string `yaml:"image,omitempty"`

	// 以下内容不存在时，则会使用全局的默认选项
	Authors  []*Author `yaml:"author,omitempty"`
	License  *Link     `yaml:"license,omitempty"`
	Template string    `yaml:"template,omitempty"`
	Language string    `yaml:"language,omitempty"`

	Content   string `yaml:"-"` // markdown 内容
	Next      *Post  `yaml:"-"`
	Prev      *Post  `yaml:"-"`
	Permalink string `yaml:"-"`
	Slug      string `yaml:"-"`
}

func (data *Data) loadPosts() error {
	paths := make([]string, 0, 10)
	err := filepath.Walk(data.Dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.ToLower(filepath.Ext(info.Name())) == ".md" {
			paths = append(paths, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	if len(paths) == 0 {
		return nil
	}

	data.Posts = make([]*Post, 0, len(paths))

	for _, path := range paths {
		post, err := data.loadPost(data.Dir, path)
		if err != nil {
			return err
		}
		data.Posts = append(data.Posts, post)
	}

	checkPosts(data.Posts)

	return nil
}

func (data *Data) loadPost(dir, path string) (*Post, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ctx := parser.NewContext()
	buf := new(bytes.Buffer)
	if err := markdown.Convert(bs, buf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	metadata, err := yaml.Marshal(meta.Get(ctx))
	if err != nil {
		return nil, err
	}
	post := &Post{}
	if err := yaml.Unmarshal(metadata, post); err != nil {
		return nil, err
	}
	post.Content = buf.String()

	if err := post.sanitize(dir, path, data.Config); err != nil {
		err.File = path
		return nil, err
	}
	return post, nil
}

func (p *Post) sanitize(dir, path string, conf *Config) *FieldError {
	if p.Title == "" {
		return &FieldError{Field: "title", Message: "不能为空"}
	}
	p.HTMLTitle = p.Title + conf.titleSuffix

	slug := strings.TrimPrefix(path, dir)
	if len(slug) > 3 && strings.ToLower(slug[len(slug)-3:]) == ".md" {
		slug = slug[:len(slug)-3]
	}
	slug = strings.Trim(filepath.ToSlash(slug), "./")
	if strings.IndexFunc(slug, func(r rune) bool { return unicode.IsSpace(r) }) >= 0 {
		return &FieldError{Field: "slug", Message: "不能包含空格"}
	}
	p.Slug = slug
	p.Permalink = conf.BuildURL(slug + ".xml")

	if len(p.TagString) == 0 {
		return &FieldError{Field: "tags", Message: "不能为空"}
	}

	// state
	if p.State != StateDefault && p.State != StateLast && p.State != StateTop {
		return &FieldError{Message: "无效的值", Field: "state"}
	}

	for i, a := range p.Authors {
		if err := a.sanitize(); err != nil {
			err.Field = "author[" + strconv.Itoa(i) + "]." + err.Field
			return err
		}
	}

	if p.License != nil {
		if err := p.License.sanitize(); err != nil {
			err.Field = "license." + err.Field
			return err
		}
	}

	return nil
}

// 对文章进行排序，需保证 created 已经被初始化
func checkPosts(posts []*Post) *FieldError {
	sort.SliceStable(posts, func(i, j int) bool {
		switch {
		case (posts[i].State == StateTop) || (posts[j].State == StateLast):
			return true
		case (posts[i].State == StateLast) || (posts[j].State == StateTop):
			return false
		default:
			return posts[i].Created.After(posts[j].Created)
		}
	})

	for _, p := range posts {
		cnt := sliceutil.Count(posts, func(i int) bool {
			return p.Slug == posts[i].Slug && p.Slug != posts[i].Slug
		})
		if cnt > 1 {
			return &FieldError{Message: "存在重复的值", Field: "slug"}
		}
	}

	// 生成 prev 和 next
	max := len(posts)
	for i := 0; i < max; i++ {
		post := posts[i]
		if i > 0 {
			post.Prev = posts[i-1]
		}
		if i < max-1 {
			post.Next = posts[i+1]
		}
	}

	return nil
}
