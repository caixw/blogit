#!/bin/sh

# SPDX-FileCopyrightText: 2020-2024 caixw
# SPDX-License-Identifier: MIT

cd $(dirname $0) || exit
date=$(date -u '+%Y%m%d')
hash=$(git rev-parse HEAD)
go build -ldflags "-X github.com/caixw/blogit/v2/internal/vars.metadata=${date}.${hash}" -v -o ./cmd/blogit/blogit ./cmd/blogit
