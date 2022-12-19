// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2/internal/testdata"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t, false)

	posts, err := LoadPosts(testdata.Source, false)
	a.NotError(err).Equal(3, len(posts))

	posts, err = LoadPosts(testdata.Source, true)
	a.NotError(err).Equal(4, len(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t, false)

	post, err := loadPost(testdata.Source, "posts/2020/12/p3.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p3").Equal(post.Slug, "posts/2020/12/p3")

	post, err = loadPost(testdata.Source, "posts/p1.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p1").Equal(post.Slug, "posts/p1").Equal(post.JSONLD, `{
    "@context": "https://schema.org/"
}
`)
}

func TestSlug(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(Slug("posts/p1.md"), "posts/p1.md")

	a.Panic(func() {
		Slug("posts/../p1.md")
	})

	a.Panic(func() {
		Slug("./posts/../p1.md")
	})

	a.Panic(func() {
		Slug("posts/")
	})
}
