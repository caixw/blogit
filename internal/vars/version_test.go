// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package vars

import (
	"testing"

	"github.com/issue9/assert/v4"
	v "github.com/issue9/version"
)

func TestVersion(t *testing.T) {
	a := assert.New(t, false)

	a.True(v.SemVerValid(Version()))
	a.True(v.SemVerValid(FullVersion()))
}
