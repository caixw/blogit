// SPDX-License-Identifier: MIT

package builder

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/filesystem"
)

var (
	_ WritableFS = MemoryFS()
	_ WritableFS = DirFS("./")
)

func testWritableFS(wfs WritableFS, a *assert.Assertion) {
	// 文件不存在
	file, err := wfs.Open("dir1/dir2/file.png")
	a.ErrorIs(err, os.ErrNotExist).Nil(file)

	// 写入文件，父目录不存在，内容为空
	a.NotError(wfs.WriteFile("dir1/dir2/file.png", nil, fs.ModePerm))
	file, err = wfs.Open("dir1/dir2/file.png")
	a.NotError(err).NotNil(file)
	data, err := io.ReadAll(file)
	a.NotError(err).Empty(data)
	a.NotError(file.Close())

	// 写入文件，父目录存在，内容不为空
	a.NotError(wfs.WriteFile("dir1/file.png", []byte{1, 2, 3}, fs.ModePerm))
	file, err = wfs.Open("dir1/file.png")
	a.NotError(err).NotNil(file)
	data, err = io.ReadAll(file)
	a.NotError(err).Equal(data, []byte{1, 2, 3})
	a.NotError(file.Close())

	// 重置后内容不再存在
	a.NotError(wfs.Reset())
	a.False(filesystem.Exists(wfs, "dir1/file.png"))
	a.False(filesystem.Exists(wfs, "dir1/dir2/file.png"))
	a.NotError(wfs.WriteFile("dir1/file.png", []byte{1, 2, 3}, fs.ModePerm)) // 重新写入
	file, err = wfs.Open("dir1/file.png")
	a.NotError(err).NotNil(file)
	a.NotError(file.Close())

	// 无效的 path 参数
	_, err = wfs.Open("/dir1/file.png")
	a.Error(err)
	err = wfs.WriteFile("/dir1/file.png", []byte{1, 2, 3}, fs.ModePerm)
	a.Error(err)
}

func TestWritableFS(t *testing.T) {
	a := assert.New(t)

	testWritableFS(MemoryFS(), a)

	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)
	testWritableFS(DirFS(dir), a)
}

func TestDir(t *testing.T) {
	a := assert.New(t)

	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)

	// 创建两个已经存在的文件
	git := path.Join(dir, ".git")
	a.NotError(os.Mkdir(git, os.ModePerm))
	obj := path.Join(git, "obj")
	a.NotError(os.WriteFile(obj, []byte{1, 2, 3}, os.ModePerm))

	inst := DirFS(dir)
	a.NotNil(inst)
	a.NotError(inst.WriteFile("a.txt", []byte{1, 2, 3}, os.ModePerm))
	a.NotError(inst.Reset())
	a.False(filesystem.Exists(inst, "a.txt"))

	// 非 writeFile 创建的文件依然存在
	_, err = os.Stat(obj)
	a.True(err == nil || errors.Is(err, fs.ErrExist))
}
