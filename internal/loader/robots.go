// SPDX-License-Identifier: MIT

package loader

import "github.com/issue9/localeutil"

// Agent 表示 robots.txt 每个代理项的内容
type Agent struct {
	Agent    []string `yaml:"agent"`
	Disallow []string `yaml:"disallow,omitempty"`
	Allow    []string `yaml:"allow,omitempty"`
}

func (a *Agent) sanitize() *FieldError {
	if len(a.Agent) == 0 {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "agent"}
	}

	if len(a.Disallow) == 0 && len(a.Allow) == 0 {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "disallow"}
	}

	return nil
}
