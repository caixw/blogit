// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestSitemap_sanitize(t *testing.T) {
	a := assert.New(t, false)

	s := &Sitemap{}
	a.Error(s.sanitize())

	s.Title = "sitemap"
	s.Priority = -1.0
	a.Error(s.sanitize())
	s.Priority = 1.1
	a.Error(s.sanitize())

	s.Priority = .8
	s.PostPriority = 0.9
	s.Changefreq = "never"
	s.PostChangefreq = "never"
	a.NotError(s.sanitize())

	s.PostChangefreq = "not-exists"
	err := s.sanitize()
	a.Equal(err.Field, "postChangefreq")
}
