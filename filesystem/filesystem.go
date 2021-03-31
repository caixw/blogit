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

	// 将内容写入文件
	//
	// path 遵守 fs.FS.Open 中有关 path 参数的处理规则。
	// 整个函数处理逻辑应该与 os.WriteFile 相同。
	// 如果文件父目录不存在，应该要自动创建。
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// 重置系统
	//
	// 该操作会删除通过 WriteFile 添加的文件。
	// 原来已经存在的内容，则依然会存在，比如 Dir() 创建的实例，
	// 就有可能是基于一个已经存在的非空目录，Reset 不应该破坏其原来的结构。
	Reset() error
}

// Memory 返回以内存作为保存对象的文件系统
func Memory() WritableFS {
	return &memoryFS{FS: memfs.New()}
}

// Dir 返回以普通目录作为保存对象的文件系统
func Dir(dir string) WritableFS {
	return &dirFS{
		dir:   dir,
		files: make([]string, 0, 10),
	}
}

type memoryFS struct {
	*memfs.FS
}

// 文件系统可能是创建在一个非目录
type dirFS struct {
	dir   string
	files []string
}

func (dir *dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	return os.Open(string(dir.dir) + "/" + name)
}

func (dir *dirFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "close", Path: name, Err: os.ErrInvalid}
	}

	p := string(dir.dir) + "/" + name
	if err := os.MkdirAll(path.Dir(p), perm); err != nil {
		return err
	}

	if err := os.WriteFile(p, data, perm); err != nil {
		return err
	}

	dir.files = append(dir.files, p) // 文件写入成功之后，记录添加的文件

	return nil
}

func (dir *dirFS) Reset() error {
	for _, f := range dir.files {
		if err := os.RemoveAll(f); err != nil {
			return err
		}
	}
	return nil
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
