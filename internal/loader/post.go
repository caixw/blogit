// SPDX-License-Identifier: MIT

package loader

import (
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/v2/internal/vars"
)

// 表示 Post.State 的各类值
const (
	StateTop     = "top"   // 置顶
	StateLast    = "last"  // 放在尾部
	StateDraft   = "draft" // 表示为草稿，不会加载此条数据
	StateDefault = ""      // 默认值
)

// Post 表示文章的信息
type Post struct {
	Title    string    `yaml:"title"`
	Created  time.Time `yaml:"created"`           // 创建时间
	Modified time.Time `yaml:"modified"`          // 修改时间
	Summary  string    `yaml:"summary,omitempty"` // 摘要，同时也作为 meta.description 的内容

	// 关联的标签列表
	//
	// 标签名为各个标签的 slug 值，可以保证其唯一。
	Tags []string `yaml:"tags"`

	// State 表示文章的状态，有以下四种值：
	// - top 表示文章被置顶；
	// - last 表示文章会被放置在最后；
	// - draft 表示这是一篇草稿；
	// - 空值 按默认的方式进行处理。
	State string `yaml:"state,omitempty"`

	// 封面地址，可以为空。
	Image string `yaml:"image,omitempty"`

	// 自定义 JSON-LD 数据
	//
	// 不需要包含 <script> 标签，只需要返回 JSON 格式数据好可。
	// 如果为空，则支自己生成 BlogPosting 类型的数据。
	JSONLD string `yaml:"jsonld,omitempty"`

	// 以下内容不存在时，则会使用全局的默认选项
	Authors  []*Author `yaml:"author,omitempty"`
	License  *Link     `yaml:"license,omitempty"`
	Template string    `yaml:"template,omitempty"`
	Language string    `yaml:"language,omitempty"`
	Keywords string    `yaml:"keywords,omitempty"`

	Content string   `yaml:"-"` // markdown 内容
	Slug    string   `yaml:"-"`
	TOC     []Header `yaml:"-"`
}

// Header TOC 的每一项内容
type Header struct {
	Level  int
	Indent int
	Text   string
	ID     string
}

// LoadPosts 加载所有的文章
//
// preview 模式下会加载草稿。
func LoadPosts(f fs.FS, preview bool) ([]*Post, error) {
	paths := make([]string, 0, 10)
	err := fs.WalkDir(f, vars.PostsDir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.ToLower(path.Ext(p)) == vars.MarkdownExt {
			paths = append(paths, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, nil
	}

	posts := make([]*Post, 0, len(paths))

	for _, p := range paths {
		post, err := loadPost(f, p)
		if err != nil {
			return nil, err
		}
		if preview || post.State != StateDraft {
			posts = append(posts, post)
		}
	}

	for _, p := range posts {
		cnt := sliceutil.Count(posts, func(i int) bool {
			return p.Slug == posts[i].Slug && p.Slug != posts[i].Slug
		})
		if cnt > 1 {
			return nil, &FieldError{Message: localeutil.Phrase("duplicate value"), Field: "slug"}
		}
	}

	return posts, nil
}

func loadPost(f fs.FS, path string) (*Post, error) {
	post, err := convert(f, path)
	if err != nil {
		return nil, err
	}

	if err := post.sanitize(path); err != nil {
		err.File = path
		return nil, err
	}
	return post, nil
}

func (p *Post) sanitize(path string) *FieldError {
	if p.Title == "" {
		return &FieldError{Field: "title", Message: localeutil.Phrase("can not be empty")}
	}
	if p.State == StateDraft { // 对草稿稍微做一下标记
		p.Title = "**" + p.Title + "**"
	}

	slug := Slug(path)
	if strings.HasSuffix(strings.ToLower(slug[len(slug)-3:]), vars.MarkdownExt) {
		slug = slug[:len(slug)-len(vars.MarkdownExt)] // 不能用 strings.TrimSuffix，后缀名可能是大写的
	}
	if strings.IndexFunc(slug, func(r rune) bool { return unicode.IsSpace(r) }) >= 0 {
		return &FieldError{Field: "slug", Message: localeutil.Phrase("can not contain spaces"), Value: slug}
	}
	if !strings.HasPrefix(slug, vars.PostsDir+"/") {
		return &FieldError{Field: "slug", Message: localeutil.Phrase("post must in", vars.PostsDir), Value: slug}
	}
	p.Slug = slug

	if len(p.Tags) == 0 {
		return &FieldError{Field: "tags", Message: localeutil.Phrase("can not be empty")}
	}

	// state
	if p.State != StateDefault && p.State != StateLast && p.State != StateTop && p.State != StateDraft {
		return &FieldError{Message: localeutil.Phrase("invalid value"), Field: "state", Value: p.State}
	}

	// template
	if p.Template == "" {
		p.Template = vars.DefaultTemplate
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

// Slug 根据文章路径返回文章的唯一 ID
func Slug(p string) string {
	if !fs.ValidPath(p) {
		panic(fmt.Sprintf("无效的参数 p: %s", p))
	}
	return strings.TrimLeft(p, "./")
}
