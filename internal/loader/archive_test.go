// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestArchive_sanitize(t *testing.T) {
	a := assert.New(t)

	archive := &Archive{}
	a.Error(archive.sanitize())

	archive = &Archive{Title: "archive"}
	a.NotError(archive.sanitize())
	a.Equal(archive.Order, OrderDesc).
		Equal(archive.Type, ArchiveTypeYear).
		Equal(archive.Format, defaultArchiveFormats[archive.Type])

	archive = &Archive{Title: "archive", Type: "not-exists"}
	a.Equal(archive.sanitize().Field, "type")

	archive = &Archive{Title: "archive", Type: ArchiveTypeMonth, Order: "not-exists"}
	a.Equal(archive.sanitize().Field, "order")

	archive = &Archive{Title: "archive", Type: ArchiveTypeMonth}
	a.NotError(archive.sanitize()).
		Equal(archive.Format, defaultArchiveFormats[archive.Type])
}
