// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)
	fs := os.DirFS("../../testdata/src")

	posts, err := LoadPosts(fs)
	a.NotError(err).Equal(3, len(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)
	fs := os.DirFS("../../testdata/src")

	post, err := loadPost(fs, "posts/2020/12/p3.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p3").Equal(post.Slug, "posts/2020/12/p3")

	post, err = loadPost(fs, "posts/p1.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p1").Equal(post.Slug, "posts/p1").Equal(post.JSONLD, `{
    "@context": "https://schema.org/"
}
`)
}

func TestSlug(t *testing.T) {
	a := assert.New(t)

	a.Equal(Slug("./posts/p1.md"), "posts/p1.md")
	a.Equal(Slug("./posts/p1.md"), "posts/p1.md")
}
