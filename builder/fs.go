// SPDX-License-Identifier: MIT

package builder

import (
	"io/fs"
	"os"
	"path"

	"github.com/psanford/memfs"
)

// WritableFS 带有写入功能的文件系统
//
// https://github.com/golang/go/issues/45757
type WritableFS interface {
	fs.FS

	// WriteFile 将内容写入文件
	//
	// path 遵守 fs.FS.Open 中有关 path 参数的处理规则。
	// 整个函数处理逻辑应该与 os.WriteFile 相同。
	// 如果文件父目录不存在，应该要自动创建。
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// Reset 重置内容
	//
	// 该操作会删除通过 WriteFile 添加的文件。
	// 原来已经存在的内容，则依然会存在，比如 Dir() 创建的实例，
	// 就有可能是基于一个已经存在的非空目录，Reset 不应该破坏其原来的结构。
	Reset() error
}

// MemoryFS 返回以内存作为保存对象的文件系统
func MemoryFS() WritableFS { return &memoryFS{FS: memfs.New()} }

// DirFS 返回以普通目录作为保存对象的文件系统
func DirFS(dir string) WritableFS {
	return &dirFS{
		FS:    os.DirFS(dir),
		dir:   dir,
		files: make([]string, 0, 10),
	}
}

// 内存文件系统，每次创建都是新的，不存在与 dirFS 一样的问题。
type memoryFS struct {
	*memfs.FS
}

type dirFS struct {
	fs.FS
	dir   string
	files []string
}

func (dir *dirFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "close", Path: name, Err: fs.ErrInvalid}
	}

	p := path.Join(dir.dir, name)
	if err := os.MkdirAll(path.Dir(p), perm); err != nil {
		return err
	}

	if err := os.WriteFile(p, data, perm); err != nil {
		return err
	}

	// BUG(caixw): 仅记录了文件，但是并未记录文件的父目录结构。
	// 在 Reset 删除时，可能会留下一堆空目录。
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
