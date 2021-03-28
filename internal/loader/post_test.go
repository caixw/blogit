// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)

	posts, err := LoadPosts("../../testdata/src")
	a.NotError(err).Equal(3, len(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	post, err := loadPost("../../testdata/src", "../../testdata/src/posts/2020/12/p3.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p3").Equal(post.Slug, "posts/2020/12/p3")

	post, err = loadPost("../../testdata/src", "../../testdata/src/posts/p1.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p1").Equal(post.Slug, "posts/p1").Equal(post.JSONLD, `{
    "@context": "https://schema.org/"
}
`)
}

func TestSlug(t *testing.T) {
	a := assert.New(t)

	a.Equal(Slug("./", "./posts/p1.md"), "posts/p1.md")
	a.Equal(Slug(".", "./posts/p1.md"), "posts/p1.md")
	a.Equal(Slug("./posts", "./posts/p1.md"), "p1.md")
	a.Equal(Slug("posts", "./posts/p1.md"), "p1.md")
	a.Equal(Slug("p", "./posts/p1.md"), "posts/p1.md")
}
