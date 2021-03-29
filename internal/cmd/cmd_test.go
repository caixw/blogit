// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"

	"github.com/issue9/assert"
)

func TestGetWD(t *testing.T) {
	a := assert.New(t)

	wfs, err := getWD()
	a.NotError(err).NotNil(wfs)
}
