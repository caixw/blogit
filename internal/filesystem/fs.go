// SPDX-License-Identifier: MIT

// Package filesystem 提供文件系统相关功能
package filesystem

import (
	"io/fs"
	"os"
	"path"

	"github.com/psanford/memfs"
)

// WritableFS 带有写入功能的文件系统
type WritableFS interface {
	fs.FS
	WriteFile(path string, data []byte, perm fs.FileMode) error
}

func Memory() WritableFS {
	return &mem{FS: memfs.New()}
}

func Dir(dir string) WritableFS {
	return dirFS(dir)
}

type mem struct {
	*memfs.FS
}

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	return os.Open(string(dir) + "/" + name)
}

func (dir dirFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	return os.WriteFile(string(dir)+"/"+name, data, perm)
}

// WriteFile
func (m *mem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	dir := path.Dir(name)
	if err := m.FS.MkdirAll(dir, perm); err != nil {
		return err
	}

	return m.FS.WriteFile(name, data, perm)
}
