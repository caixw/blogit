// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"testing"

	"github.com/issue9/assert"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/loader"
)

func TestWriteYAML(t *testing.T) {
	a := assert.New(t)

	obj := &loader.Theme{Description: "desc", URL: "https://example.com"}

	wfs := filesystem.Memory()
	a.NotError(writeYAML(wfs, "conf.yaml", obj))

	f, err := wfs.Open("conf.yaml")
	a.NotError(err).NotNil(f)
	data, err := io.ReadAll(f)
	a.NotError(err).NotNil(data)

	inst := &loader.Theme{}
	a.NotError(yaml.Unmarshal(data, inst))
	a.Equal(inst, obj)
}
