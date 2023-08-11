// SPDX-License-Identifier: MIT

package loader

// Agent 表示 robots.txt 每个代理项的内容
type Agent struct {
	Agent    []string `yaml:"agent"`
	Disallow []string `yaml:"disallow,omitempty"`
	Allow    []string `yaml:"allow,omitempty"`
}

func (a *Agent) sanitize() *FieldError {
	if len(a.Agent) == 0 {
		return &FieldError{Message: Required, Field: "agent"}
	}

	if len(a.Disallow) == 0 && len(a.Allow) == 0 {
		return &FieldError{Message: Required, Field: "disallow"}
	}

	return nil
}
