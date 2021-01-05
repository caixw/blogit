// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)

	posts, err := LoadPosts("../testdata")
	a.NotError(err).Equal(3, len(posts))

	posts, err = LoadPosts("../testdata/posts/2020/12")
	a.NotError(err).Equal(1, len(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	post, err := loadPost("../testdata", "../testdata/posts/2020/12/p3.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p3").Equal(post.Slug, "posts/2020/12/p3")
}
