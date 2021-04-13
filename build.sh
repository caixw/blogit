#!/bin/sh

# SPDX-License-Identifier: MIT

cd $(dirname $0) || exit
date=$(date -u '+%Y%m%d')
hash=$(git rev-parse HEAD)
go build -ldflags "-X github.com/caixw/blogit/internal/vars.metadata=${date}.${hash}" -v -o ./cmd/blogit/blogit ./cmd/blogit
