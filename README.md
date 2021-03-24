# blogit

[![Test](https://github.com/caixw/blogit/workflows/Test/badge.svg)](https://github.com/caixw/blogit/actions?query=workflow%3ATest)
[![Go version](https://img.shields.io/github/go-mod/go-version/caixw/blogit)](https://golang.org)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/caixw/blogit)](https://pkg.go.dev/github.com/caixw/blogit)

静态博客生成工具

## 安装

常用平台可以从 <https://github.com/caixw/blogit/releases> 下载，并将二进制文件放入 `PATH` 即可。

如果不存在你当前平台的二进制，可以自己编译：

```shell
git clone https://github.com/caixw/blogit.git
cd blogit
./build.sh
```

## 使用

`testdata/src/` 下包含了一个完整的数据源，你可以直接复制该目录下的内容稍作修改，
即可以当作自己的博客数据源。

## docker

目前 docker 同时托管于 [docker.io](https://hub.docker.com/r/caixw/blogit) 和 [ghcr.io](https://ghcr.io/caixw/blogit)，可通过以下方式获取相应在的容器：

`docker pull docker.io/caixw/blogit:latest`

`docker pull ghcr.io/caixw/blogit:latest`

## 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
