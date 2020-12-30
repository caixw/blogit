// SPDX-License-Identifier: MIT

package builder

import (
	"sort"
	"time"

	"github.com/caixw/blogit/internal/loader"
)

// Archive 表示某一时间段的存档信息
type Archive struct {
	date  time.Time // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string    // 当前存档页的标题
	Posts []*Post   // 当前存档的文章列表
}

func buildArchives(conf *loader.Config, posts []*Post) ([]*Archive, error) {
	archives := make([]*Archive, 0, 10)

	for _, post := range posts {
		t := post.Created
		var date time.Time

		switch conf.Archive.Type {
		case loader.ArchiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location())
		case loader.ArchiveTypeYear:
			date = time.Date(t.Year(), 2, 0, 0, 0, 0, 0, t.Location())
		default:
			panic("无效的 archive.type 值")
		}

		found := false
		for _, archive := range archives {
			if archive.date.Equal(date) {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			archives = append(archives, &Archive{
				date:  date,
				Title: date.Format(conf.Archive.Format),
				Posts: []*Post{post},
			})
		}
	} // end for

	sort.SliceStable(archives, func(i, j int) bool {
		if conf.Archive.Order == loader.ArchiveOrderDesc {
			return archives[i].date.After(archives[j].date)
		}
		return archives[i].date.Before(archives[j].date)
	})

	return archives, nil
}
