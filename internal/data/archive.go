// SPDX-License-Identifier: MIT

package data

import (
	"sort"
	"time"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Archives 存档的相关信息
type Archives struct {
	Title       string
	Permalink   string
	Keywords    string
	Description string
	Archives    []*Archive
}

// Archive 单个存档项
type Archive struct {
	date  time.Time // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string    // 当前存档页的标题
	Posts []*Post   // 当前存档的文章列表
}

func buildArchives(conf *loader.Config, posts []*Post) (*Archives, error) {
	if conf.Archive == nil {
		return nil, nil
	}

	list := make([]*Archive, 0, 10)
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
		for _, archive := range list {
			if archive.date.Equal(date) {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			list = append(list, &Archive{
				date:  date,
				Title: date.Format(conf.Archive.Format),
				Posts: []*Post{post},
			})
		}
	} // end for

	sort.SliceStable(list, func(i, j int) bool {
		if conf.Archive.Order == loader.OrderDesc {
			return list[i].date.After(list[j].date)
		}
		return list[i].date.Before(list[j].date)
	})

	if conf.Archive.Keywords == "" {
		conf.Archive.Keywords = conf.Keywords
	}

	if conf.Archive.Description == "" {
		conf.Archive.Description = conf.Description
	}

	return &Archives{
		Title:       conf.Archive.Title,
		Permalink:   buildURL(conf.URL, vars.ArchiveFilename),
		Keywords:    conf.Archive.Keywords,
		Description: conf.Archive.Description,
		Archives:    list,
	}, nil
}
