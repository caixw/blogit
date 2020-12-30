// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/loader"
)

// Data 处理后的数据
type Data struct {
	Title       string
	Subtitle    string
	TitleSuffix string // 每篇文章标题的后缀
	Icon        *Icon
	Menus       []*Menu
	Theme       string

	Uptime   time.Time
	Builded  time.Time // 最后次编译时间
	Created  time.Time
	Modified time.Time

	Tags  []*Tag
	Posts []*Post
}

// Icon 图标信息
type Icon = loader.Icon

// License 表示链接信息
type License = loader.License

// Menu 采单项
type Menu = loader.Menu

// Author 表示作者信息
type Author = loader.Author
