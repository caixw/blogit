// SPDX-License-Identifier: MIT

package blogit

import (
	"log"
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_init(t *testing.T) {
	a := assert.New(t)

	// 都采用默认值
	o := &Options{}
	a.False(o.inited).NotError(o.init()).True(o.inited).
		Equal(o.Path, "/").
		Equal(o.Src, "./").
		NotNil(o.Dest).NotNil(o.b).
		Equal(o.Addr, ":http")
	a.NotError(o.init()).NotError(o.init())

	o = &Options{
		BaseURL: "https://localhost:8080/path/",
	}
	a.NotError(o.init()).
		Equal(o.Path, "/path/").
		Equal(o.Addr, ":8080")

	o = &Options{
		BaseURL: "https://localhost/path/",
	}
	a.NotError(o.init()).
		Equal(o.Path, "/path/").
		Equal(o.Addr, ":443")

	o = &Options{
		BaseURL: "http://localhost/path/",
		Cert:    "./cert",
		Key:     "./key",
	}
	a.NotError(o.init()).
		Equal(o.Path, "/path/").
		Equal(o.Addr, ":http")

	o = &Options{
		Cert: "./cert",
		Key:  "./key",
		Info: log.Default(),
		Erro: log.Default(),
		Succ: log.Default(),
	}
	a.NotError(o.init()).
		Equal(o.Path, "/").
		Equal(o.Addr, ":443")

	o = &Options{
		BaseURL: "ftp://localhost/path/",
	}
	a.ErrorString(o.init(), "不支持")

	// url 格式错误
	o = &Options{
		BaseURL: "http://localh%2%ost/path/",
	}
	a.Error(o.init())
}
