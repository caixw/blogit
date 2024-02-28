// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package locales 提供本地化相关操作
package locales

import "embed"

//go:embed *.yaml
var locales embed.FS

func Locales() embed.FS { return locales }
