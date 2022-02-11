// SPDX-License-Identifier: MIT

package console

import (
	"github.com/caixw/blogit/v2/locale"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"
)

func NewPrinter() (*message.Printer, error) {
	b := catalog.NewBuilder()
	if err := localeutil.LoadMessageFromFSGlob(b, locale.Locales(), "*.yaml", yaml.Unmarshal); err != nil {
		return nil, err
	}

	systag, _ := localeutil.DetectUserLanguageTag() // 即使出错，依然会返回 language.Tag
	return message.NewPrinter(systag, message.Catalog(b)), nil
}
