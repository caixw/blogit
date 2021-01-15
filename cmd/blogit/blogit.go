// SPDX-License-Identifier: MIT

package main

import "github.com/caixw/blogit/internal/cmd"

func main() {
	if err := cmd.Exec(); err != nil {
		panic(err)
	}
}
