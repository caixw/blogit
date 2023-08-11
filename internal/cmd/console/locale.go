// SPDX-License-Identifier: MIT

package console

import (
	"github.com/issue9/localeutil/message/serialize"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v3"

	"github.com/caixw/blogit/v2/locales"
)

func NewPrinter(tag language.Tag) (*message.Printer, error) {
	b := catalog.NewBuilder()
	l, err := serialize.LoadFSGlob(locales.Locales(), "*.yaml", yaml.Unmarshal)
	if err != nil {
		return nil, err
	}
	for _, ll := range l {
		if err := ll.Catalog(b); err != nil {
			return nil, err
		}
	}

	return message.NewPrinter(tag, message.Catalog(b)), nil
}
