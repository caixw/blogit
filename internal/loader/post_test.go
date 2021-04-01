// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/testdata"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)

	posts, err := LoadPosts(testdata.Source)
	a.NotError(err).Equal(3, len(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

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
	a := assert.New(t)

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
