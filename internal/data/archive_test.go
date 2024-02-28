// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package data

import (
	"testing"
	"time"

	"github.com/issue9/assert/v4"

	"github.com/caixw/blogit/v2/internal/loader"
)

func TestBuildArchives(t *testing.T) {
	a := assert.New(t, false)
	const month = 30 * 24 * time.Hour

	as, err := buildArchives(&loader.Config{}, nil)
	a.NotError(err).Nil(as)

	dt := time.Date(2020, 03, 1, 0, 0, 0, 0, time.Local) // 3 月份，不会跨年
	posts := []*Post{
		{
			Title:   "1",
			Created: dt,
		},
		{
			Title:   "2",
			Created: dt,
		},
		{
			Title:   "3",
			Created: dt.Add(-2 * month),
		},
	}

	// 按月分类
	conf := &loader.Config{
		Keywords:    "keys",
		Description: "des",
		Archive: &loader.Archive{
			Title:  "title",
			Order:  loader.OrderAsc,
			Type:   loader.ArchiveTypeMonth,
			Format: "2006-01",
		},
	}
	as, err = buildArchives(conf, posts)
	a.NotError(err).NotNil(as)
	a.Equal(as.Title, "title").Equal(as.Keywords, conf.Keywords).Equal(as.Description, conf.Description)
	a.Equal(2, len(as.Archives))                            // 按月可分成了两部分
	a.True(as.Archives[0].date.Before(as.Archives[1].date)) // 排序是否正常

	conf.Archive.Order = loader.OrderDesc
	as, err = buildArchives(conf, posts)
	a.NotError(err).NotNil(as)
	a.Equal(as.Title, "title").Equal(as.Keywords, conf.Keywords).Equal(as.Description, conf.Description)
	a.Equal(2, len(as.Archives))                           // 按月可分成了两部分
	a.True(as.Archives[0].date.After(as.Archives[1].date)) // 排序是否正常

	// 按年分类，不跨年

	conf.Archive.Type = loader.ArchiveTypeYear
	as, err = buildArchives(conf, posts)
	a.NotError(err).NotNil(as)
	a.Equal(1, len(as.Archives))

	// 文章始终是按时间顺序的
	ps := as.Archives[0].Posts
	a.Equal(ps[0].Title, "1").Equal(ps[1].Title, "2").Equal(ps[2].Title, "3")

	// 文章始终是按时间顺序的，即使改变了 Archive.Order
	conf.Archive.Type = loader.ArchiveTypeYear
	conf.Archive.Order = loader.OrderDesc
	as, err = buildArchives(conf, posts)
	a.NotError(err).NotNil(as)
	a.Equal(1, len(as.Archives))
	ps = as.Archives[0].Posts
	a.Equal(ps[0].Title, "1").Equal(ps[1].Title, "2").Equal(ps[2].Title, "3")

	// 按年份，跨年

	dt = time.Date(2020, 01, 15, 0, 0, 0, 0, time.Local)
	posts[0].Created = dt
	posts[1].Created = dt
	posts[2].Created = dt.Add(-2 * month)
	conf.Archive.Type = loader.ArchiveTypeYear
	as, err = buildArchives(conf, posts)
	a.NotError(err).NotNil(as)
	a.Equal(2, len(as.Archives))
}
