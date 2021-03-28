// SPDX-License-Identifier: MIT

package filesystem

import (
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ WritableFS = Memory()
	_ WritableFS = Dir("./")
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

	// 写入文件，父目录存在，内容不为空
	a.NotError(wfs.WriteFile("dir1/file.png", []byte{1, 2, 3}, fs.ModePerm))
	file, err = wfs.Open("dir1/file.png")
	a.NotError(err).NotNil(file)
	data, err = io.ReadAll(file)
	a.NotError(err).Equal(data, []byte{1, 2, 3})

	// 无效的 path 参数
	_, err = wfs.Open("/dir1/file.png")
	a.Error(err)
	err = wfs.WriteFile("/dir1/file.png", []byte{1, 2, 3}, fs.ModePerm)
	a.Error(err)
}

func TestWritableFS(t *testing.T) {
	a := assert.New(t)

	testWritableFS(Memory(), a)

	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)
	testWritableFS(Dir(dir), a)
}
