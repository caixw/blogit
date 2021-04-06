# 使用

当前文档描述了针对 blogit 的配置内容：

## 初始化

可通过 `blogit init path/to/blog/dir` 进行初始经，会在 `path/to/blog/dir` 生成完整的项目文件。

### conf.yaml

项目的配置文件，大部分项目的全局修改均由该文件配置。

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 网站标题
| subtitle        | string      | 网站副标题
| titleSeparator  | string      | 标题分隔符，在 html>head>title 标签中，非首页通过该符号在标题后加网站标题
| url             | string      | 网站根目录，比如 `https://example.com/blog`
| language        | string      | 网站的默认语言，在文章中若没有专门配置，则采用此值作为文章的默认语言，比如 `cmn-Hans`
| uptime          | string      | 网站的上线时间，rfc3999 格式
| icon            | Icon        | favicon 图标定义
| author          | Author      | 网站的默主作者
| license         | Link        | 网站的默认版权信息
| theme           | string      | 网站采用的主题，该名称必须是 themes/ 下的文件夹名称。
| keywords        | string      | 首页的 html>head>meta.keywords 标签的值
| description     | string      | 首页的 html>head>meta.description 标签的值
| menus           | []Link      | 菜单栏内容
| toc             | number      | 当标题数量大于此值时，才会生成 TOC 数据，默认值为 0。
| index           | Index       | 索引页相关的设置
| archive         | Archive     | 存档页的相关定义，可以为空，表示不需要该页面。
| rss             | RSS         | RSS 的相关定义，为空表示不需要。
| atom            | RSS         | atom 的相关定义，为空表示不需要。
| Sitemap         | Sitemap     | sitemap 的相关定义，为空表示不需要。
| Robots          | []Agent     | robots.txt 文件的配置，如果为空表示不需要由项目管理 robots.txt 文件。
| Profile         | Profile     | 管理 README.md 的生成，github 中仅与账号同名的项目才会在 profile 中显示。

#### Icon

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| url             | string      | 图标地址，比如 /favicon.svg
| type            | string      | 图标的类型，比如 image/svg+xml
| sizes           | string      | 图标的大小，比如 256x256

#### Author

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| name            | string      | 作者名称
| url             | string      | 作者网站
| email           | string      | 作者的邮箱
| avatar          | string      | 头像

#### Link

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| text            | string      | 链接的文本
| url             | string      | 链接指向的地址

#### RSS

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标题，部分模板可能会引用到。
| size            | number      | 生成的文章数量

#### Index

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标题，可以带 `%d` 占位符，表示当前页的页码数字。
| size            | number      | 生成的文章数量

#### Sitemap

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标题，部分模板可能会引用到。
| enableTag       | boolean     | 是否将标签页也放入 sitemap
| postPriority    | number      | 对应文章的 priority 值
| postChangefreq  | string     | 对应文章的 changefreq 值
| priority        | number      | 其它页面的 priority 值
| changefreq      | string     | 其它页面的 changefreq 值

#### Agent

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| agent           | []string    | 指定可用的爬虫，比如 ['*']
| disallow        | []string    | 禁止抓取的目录
| allow           | []string    | 允许的目录

#### Profile

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标题，该值会自动加上 ### 字符
| footer          | string      | 页脚部分，自动加上 ##### 字符
| size            | string      | 生成该数量的文章列表

该配置生成的内容大致如下：

```md
### title

[post title](post link)
[post title](post link)
[post title](post link)

##### footer
```

### 标签

blogit 不支持文章分类，也没有一般博客的页面和文章的区别，只能通过标签对文章进行归类统计。

通过 tags.yaml 可以定义标签的相关信息：

#### tags.yaml

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标签页的标题
| keywords        | string      | 标签页的 html>head>meta.keywords 中的值，如果没有，则提取所有标签值作为该值
| description     | string      | 标签页的 html>head>meta.description 元素的值。
| order           | string      | 排序，可是 `asc` 和 `desc`
| orderType       | string      | 排序方式，可以是 size 表示按关联文章数量进行排序，或是为空，按添加顺序。
| tags            | []Tag       | 标签列表

#### Tag

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 标签名称
| slug            | string      | 标签的唯一 ID，一般会显示在 URL 中
| content         | string      | 标签的描述，可以是 markdown 格式。

### 主题

主题包含在 themes/ 目录下，每个目录为一个主题，每个主题包含 `theme.yaml` 文件，
该文件配置了一主题相关的一些内容。

#### theme.yaml

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| url             | string      | 主题关联的地址
| description     | string      | 详细说明
| authors         | []Author    | 主题作者
| screenshots     | []string    | 截图
| templates       | []string    | 文章详情页的模板名称，在文章配置项指定的模板必须上此列表中的某个值。
| highlights      | Highlight   | 语法高亮的相关配置
| sitemap         | string      | 为 sitemap.xml 指定一个 xsl 转换文件
| atom            | string      | 为 atom.xml 指定一个 xsl 转换文件
| rss             | string      | 为 rss.xml 指定一个 xsl 转换文件

#### Highlight

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| name            | string      | 名称，该值可通过 `blogit styles` 获取
| media           | string      | 对应的媒体查询值，比如 `print`、`(prefers-color-scheme: dark)` 等。

#### 模板

每个主题下的 layout 目录之下可以存放模板。模板采用 go 的 html/template 语法。传递的变量如下：

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| Type            | string      | 当前页面类型，同时也是页面采用的模板名称。除了文章详情之外，其它值都是固定的。
| Site            | Site        | 站点的数据，不同页面该数据都是相同的。
| Title           | string      | 当前页的标题，出现在 html>head>title 元素中，会加上网站名称作为后缀。
| Permalink       | string      | 当前页的唯一链接
| Keywords        | string      | 当前页的 html>head>meta.keywords 元素中数据。
| Description     | string      | 当前页的 html>head>meta.description 元素中数据。
| Description     | string      | 当前页的 html>head>meta.description 元素中数据。
| Prev            | Link        | 前一页的链接
| Next            | Link        | 后一页的链接
| Author          | Author      | 当前页内容的作者
| License         | Link        | 当前页的版权信息
| Language        | string      | 当前页所采用的语言
| JSONLD          | string      | 当前页的 JSON-LD 数据
| Tag             | Tag         | 如果当前页是 `tag`，那么表示该标签的数据，否则为空值。
| Post            | Post        | 如果当前页是 `post`，那么表示该页的数据，否则为空值。
| Index           | Index       | 如果当前页是 `index`，那么表示该页的数据，否则为空值。
| Archives        | Archives    | 存档信息

Type 可以有以下值：

- tags: 标签列表
- index: 首页
- tag: 标签页
- archive: 存档页
- post: 文章详情页，文章情况页也可以是其它任意非空值。

##### Site

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| AppName         | string      | 生成该模板的程序名称
| AppVersion      | string      | 生成该模板的程序版本
| Theme           | Theme       | 主题信息
| Highlights      | []StyleLink | 代码高这所需要的 CSS 文件
| Title           | string      | 网站名称
| Subtitle        | string      | 网站副标题
| URL             | string      | 网站根地址
| Icon            | Icon        | 网站图标
| Author          | Author      | 网站的默认作者
| RSS             | Link        | RSS 链接
| Atom            | Link        | Atom 链接
| Sitemap         | Link        | Sitemap 链接
| Menus           | []Link      | 全局菜单
| Tags            | Tags        | 标签列表
| Uptime          | date        | 上线时间
| Created         | date        | 最后次创建文章的时间
| Modified        | date        | 最后次修改文章的时间
| Builded         | date        | 编译项目的时间

##### Index

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| Title           | string      | 页标题
| Permalink       | string      | 链接
| Posts           | []Post      | 文章列表
| Index           | number      | 当前页的页码
| Path            | string      | 当前页的文件地址

##### StyleLink

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| URL             | string      | CSS 指向的地址
| Media           | string      | 对应的媒体查询值，比如 `print`、`(prefers-color-scheme: dark)` 等。

## 添加文章

命令 `blogit post path/to/file` 可以在当前目录下添加 `posts/path/to/file.md` 文件，并在文件中添加必要的字段，
每个文件中可用的字段如下：

| 名称            | 类型        | 描述
|-----------------|-------------|-------------
| title           | string      | 文章标题
| created         | string      | 创建时间，rfc3339 格式
| modified        | string      | 修改时间，rfc3339 格式
| summary         | string      | 摘要
| tags            | []string    | 关联的标签
| state           | string      | 状态，可以是以下值：top 表示文章被置顶；last 表示文章会被放置在最后；draft 表示这是一篇草稿；空值 按默认的方式进行处理。
| image           | string      | 封面图片
| jsonld          | string      | 自定义 json-ld 数据，为空则会自动生成。
| authors         | []Author    | 作者，如果为空，则采用 conf.yaml 中对应的值。
| license         | string      | 文章的版权信息，如果为空则采用 conf.yaml 中对应的值。
| template        | string      | 文章的模板，如果为空，则采用默认值 `post`。
| keywords        | string      | html>head>meta.keywords 的值，如果为空，自动提取 tags 作为默认值。
| language        | string      | 页面的语言，如果为空，则采用 conf.yaml 中对应的值。
