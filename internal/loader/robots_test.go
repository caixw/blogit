// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestRobots(t *testing.T) {
	a := assert.New(t)

	agent := &Agent{}
	err := agent.sanitize()
	a.Equal(err.Field, "agent")

	agent.Agent = []string{"*"}
	err = agent.sanitize()
	a.Equal(err.Field, "disallow")
}
