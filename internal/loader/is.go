// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package loader

import (
	"net/mail"
	"net/url"
)

func isURL(u string) bool {
	_, err := url.Parse(u)
	return err == nil
}

func isEmail(u string) bool {
	_, err := mail.ParseAddress(u)
	return err == nil
}
