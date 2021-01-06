// SPDX-License-Identifier: MIT

package builder

import (
	"sort"
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
)

type archive struct {
	date  time.Time   // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string      // 当前存档页的标题
	Posts []*postMeta // 当前存档的文章列表
}

func (b *Builder) buildArchives(path string, d *data.Data) error {
	archives := make([]*archive, 0, 10)
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

		tags := make([]*tag, 0, len(post.Tags))
		for _, t := range post.Tags {
			tags = append(tags, &tag{
				Permalink: d.BuildURL("tags", t.Slug+".xml"),
				Title:     t.Title,
				Color:     t.Color,
				Content:   t.Content,
				Created:   toDatetime(t.Created, d),
				Modified:  toDatetime(t.Modified, d),
			})
		}

		pm := &postMeta{
			Permalink: d.BuildURL(post.Slug + ".xml"),
			Title:     post.Title,
			Created:   toDatetime(post.Created, d),
			Modified:  toDatetime(post.Modified, d),
			Tags:      tags,
		}

		found := false
		for _, archive := range archives {
			if archive.date.Equal(date) {
				archive.Posts = append(archive.Posts, pm)
				found = true
				break
			}
		}
		if !found {
			archives = append(archives, &archive{
				date:  date,
				Title: date.Format(d.Archive.Format),
				Posts: []*postMeta{pm},
			})
		}
	} // end for

	sort.SliceStable(archives, func(i, j int) bool {
		if d.Archive.Order == loader.ArchiveOrderDesc {
			return archives[i].date.After(archives[j].date)
		}
		return archives[i].date.Before(archives[j].date)
	})

	return b.appendXMLFile(path, d.BuildThemeURL("archive.xsl"), d.Modified, archives)
}
