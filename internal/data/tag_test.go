// SPDX-License-Identifier: MIT

package data

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/v2/internal/loader"
)

func TestSortTags(t *testing.T) {
	a := assert.New(t)

	// type=default
	tags := []*Tag{
		{
			Title: "1",
			Posts: []*Post{{}, {}, {}},
		},

		{
			Title: "2",
			Posts: []*Post{{}, {}},
		},
	}
	sortTags(tags, loader.TagOrderTypeDefault, loader.OrderAsc)
	a.Equal(tags[0].Title, "1").Equal(tags[1].Title, "2")

	sortTags(tags, loader.TagOrderTypeDefault, loader.OrderDesc)
	a.Equal(tags[0].Title, "2").Equal(tags[1].Title, "1")

	// type=size

	sortTags(tags, loader.TagOrderTypeSize, loader.OrderAsc)
	a.Equal(tags[0].Title, "1").Equal(tags[1].Title, "2")

	sortTags(tags, loader.TagOrderTypeSize, loader.OrderDesc)
	a.Equal(tags[0].Title, "2").Equal(tags[1].Title, "1")
}
