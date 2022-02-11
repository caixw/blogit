// SPDX-License-Identifier: MIT

// Package locale 提供本地化相关操作
package locale

import "embed"

//go:embed *.yaml
var locales embed.FS

func Locales() embed.FS { return locales }
