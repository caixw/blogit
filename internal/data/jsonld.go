// SPDX-License-Identifier: MIT

package data

import (
	"encoding/json"
	"time"

	"github.com/caixw/blogit/internal/loader"
)

type ldBlogPosting struct {
	ldCreativeWork
}

type ldCreativeWork struct {
	Context  string      `json:"@context"`
	Type     string      `json:"@type"`
	Headline string      `json:"headline,omitempty"`
	Name     string      `json:"name,omitempty"`
	Authors  []*ldPerson `json:"author,omitempty"`
	Created  time.Time   `json:"dateCreated,omitempty"`
	Modified time.Time   `json:"dateModified,omitempty"`
	License  string      `json:"license,omitempty"`
	Keywords string      `json:"keywords,omitempty"`
	Language string      `json:"inLanguage,omitempty"`
}

type ldPerson struct {
	Type  string `json:"@type"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

func newLDBlogPosting(p *loader.Post) *ldBlogPosting {
	blog := &ldBlogPosting{
		ldCreativeWork: ldCreativeWork{
			Context:  "https://schema.org/",
			Type:     "BlogPosting",
			Headline: p.Title,
			Created:  p.Created,
			Modified: p.Modified,
			License:  p.License.URL,
			Keywords: p.Keywords,
			Language: p.Language,
		},
	}

	for _, a := range p.Authors {
		blog.Authors = append(blog.Authors, &ldPerson{
			Type:  "Person",
			Name:  a.Name,
			Email: a.Email,
			URL:   a.URL,
		})
	}

	return blog
}

func buildPostLD(p *loader.Post) (string, error) {
	data, err := json.Marshal(newLDBlogPosting(p))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
