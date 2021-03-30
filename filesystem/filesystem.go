// SPDX-License-Identifier: MIT

// Package filesystem 提供文件系统相关功能
package filesystem

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/psanford/memfs"
)

// WritableFS 带有写入功能的文件系统
type WritableFS interface {
	fs.FS

	// 将内容写入文件
	//
	// path 遵守 fs.FS.Open 中有关 path 参数的处理规则。
	// 整个函数处理逻辑应该与 os.WriteFile 相同。
	// 如果文件父目录不存在，应该要自动创建。
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// 清空内容
	Reset() error
}

// Memory 返回以内存作为保存对象的文件系统
func Memory() WritableFS {
	return &memoryFS{FS: memfs.New()}
}

// Dir 返回以普通目录作为保存对象的文件系统
func Dir(dir string) WritableFS {
	return dirFS(dir)
}

type memoryFS struct {
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
		return &fs.PathError{Op: "close", Path: name, Err: os.ErrInvalid}
	}

	path := string(dir) + "/" + filepath.Dir(name)
	if err := os.MkdirAll(path, perm); err != nil {
		return err
	}

	return os.WriteFile(string(dir)+"/"+name, data, perm)
}

func (dir dirFS) Reset() error {
	if err := os.RemoveAll(string(dir)); err != nil {
		return err
	}

	return os.Mkdir(string(dir), os.ModePerm)
}

func (m *memoryFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if err := m.FS.MkdirAll(path.Dir(name), perm); err != nil {
		return err
	}
	return m.FS.WriteFile(name, data, perm)
}

func (m *memoryFS) Reset() error {
	m.FS = memfs.New()
	return nil
}
