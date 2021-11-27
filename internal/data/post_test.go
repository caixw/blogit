// SPDX-License-Identifier: MIT

package data

import (
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/blogit/v2/internal/loader"
)

func TestSortPosts(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()

	posts := []*loader.Post{
		{
			Title:   "1",
			State:   loader.StateLast,
			Created: now,
		},
		{
			Title:   "2",
			State:   loader.StateTop,
			Created: now,
		},
		{
			Title:   "3",
			Created: now.Add(-time.Hour),
		},
	}
	sortPosts(posts)

	a.Equal(posts[0].Title, "2").
		Equal(posts[1].Title, "3").
		Equal(posts[2].Title, "1")
}

func TestSortPostsByCreated(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()

	posts := []*Post{
		{
			Title:   "1",
			Created: now,
		},
		{
			Title:   "2",
			Created: now.Add(-time.Hour),
		},
		{
			Title:   "3",
			Created: now.Add(-time.Hour),
		},
		{
			Title:   "4",
			Created: now,
		},
	}

	ps := sortPostsByCreated(posts)
	a.Equal(ps[0].Title, "1").
		Equal(ps[1].Title, "4").
		Equal(ps[2].Title, "2").
		Equal(ps[3].Title, "3")

	// posts 不会改变
	a.Equal(posts[0].Title, "1").
		Equal(posts[1].Title, "2").
		Equal(posts[2].Title, "3").
		Equal(posts[3].Title, "4")
}
