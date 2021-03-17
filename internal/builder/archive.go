// SPDX-License-Identifier: MIT

package builder

import (
	"sort"
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

type archive struct {
	date  time.Time    // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string       // 当前存档页的标题
	Posts []*data.Post // 当前存档的文章列表
}

func (b *builder) buildArchive(path string, d *data.Data) error {
	if d.Archive == nil {
		return nil
	}

	list := make([]*archive, 0, 10)
	for _, post := range d.Posts {
		t := post.Created
		var date time.Time

		switch d.Archive.Type {
		case loader.ArchiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location())
		case loader.ArchiveTypeYear:
			date = time.Date(t.Year(), 2, 0, 0, 0, 0, 0, t.Location())
		default:
			panic("无效的 archive.type 值")
		}

		found := false
		for _, archive := range list {
			if archive.date.Equal(date) {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			list = append(list, &archive{
				date:  date,
				Title: date.Format(d.Archive.Format),
				Posts: []*data.Post{post},
			})
		}
	} // end for

	sort.SliceStable(list, func(i, j int) bool {
		if d.Archive.Order == loader.ArchiveOrderDesc {
			return list[i].date.After(list[j].date)
		}
		return list[i].date.Before(list[j].date)
	})

	p := b.page(vars.ArchiveTemplate)
	p.Title = d.Archive.Title
	p.Permalink = d.BuildURL(vars.ArchiveFilename)
	p.Keywords = "TODO"
	p.Description = "TODO"
	p.Language = d.Language

	return b.appendTemplateFile(vars.ArchiveFilename, p)
}
