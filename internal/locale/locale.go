// SPDX-License-Identifier: MIT

// Package locale 提供本地化相关操作
package locale

import (
	"embed"
	"io/fs"

	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"
)

//go:embed *.yaml
var locales embed.FS

var b *catalog.Builder

func NewPrinter() (*message.Printer, error) {
	systag, _ := localeutil.DetectUserLanguageTag() // 即使出错，依然会返回 language.Tag

	if b == nil {
		matches, err := fs.Glob(locales, "*.yaml")
		if err != nil {
			return nil, err
		}

		b = catalog.NewBuilder()

		for _, file := range matches {
			if err := localeutil.LoadMessageFromFS(b, locales, file, yaml.Unmarshal); err != nil {
				return nil, err
			}
		}
	}

	return message.NewPrinter(systag, message.Catalog(b)), nil
}
