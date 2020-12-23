// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)

	data := &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("./conf.yaml"))
	err := data.loadPosts()
	a.NotError(err).Equal(3, len(data.Posts))

	data = &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("./conf.yaml"))
	data.Dir = "./testdata/posts/2020/12"
	err = data.loadPosts()
	a.NotError(err).Equal(1, len(data.Posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	data := &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("./conf.yaml"))
	post, err := data.loadPost("./testdata", "./testdata/posts/2020/12/p3.md")
	a.NotError(err).NotNil(post)
	a.Equal(post.Title, "p3").Equal(post.Slug, "posts/2020/12/p3")
}
