// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := LoadTheme("./testdata", "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 2).Equal(theme.ID, "default")

	theme, err = LoadTheme("./testdata", "not-exists")
	a.ErrorIs(err, os.ErrNotExist).Nil(theme)
}
