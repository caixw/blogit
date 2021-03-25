<?xml version="1.0" encoding="utf-8"?>
<!--
为 sitemap.xml 产生一个比较美观的人机界面。

@author     caixw <https://caixw.io>
@license    MIT License
-->
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:sm="http://www.sitemaps.org/schemas/sitemap/0.9">
<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<xsl:template match="/">
<html lang="zh-cmn-Hans">
<head>
<title>XML Sitemap</title>
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="generator" content="https://caixw.io" />
<link rel="stylesheet" type="text/css" href="/themes/default/style.css" />
<style>
/* 颜色等值由引用的 style.css 决定。当前 style 标签中只规定了表格的相关设置。 */
table {
    width: 100%;
    text-align: left;
    border-collapse: collapse;
    line-height: 1.5;
}
td, th {
    padding: 1px 5px;
}
td.sitemap-lastmod {
    text-align:right;
}
.lastmod {
    width: 10rem;
}
.changefreq {
    width: 5rem;
}
.priority {
    width: 3rem;
}
</style>
</head>
<body id="main">
    <header>
        <h1>XML Sitemap</h1>
        <div>
        <p>这是个标准的 sitemap 文件。您可以将此文件提交给 <a href="https://www.google.com/webmasters/tools/">Google</a>、<a href="https://www.bing.com/webmaster">Bing</a> 或<a href="https://sitemap.baidu.com">百度</a>，让搜索引擎更好地收录您的网站内容。</p>

        <p>若是存在 sitemap 的索引文件，则只需提交索引文件。更详细的信息请参考<a href="https://www.sitemaps.org/zh_CN/protocol.php">这里</a>。</p>
        </div>
    </header>

    <xsl:apply-templates select="sm:urlset" />

    <footer>
        <p>此 <a href="https://github.com/caixw/sitemap.xsl">XSL 模板</a>由 <a href="https://caixw.io">caixw</a> 制作，并采用 <a href="https://www.opensource.org/licenses/MIT">MIT</a> 开源许可证发布。</p>
    </footer>
</body>
</html>
</xsl:template>


<xsl:template match="sm:urlset">

<xsl:variable name="max">
  <xsl:for-each select="/sm:urlset/sm:url/sm:lastmod">
    <xsl:sort select="." order="descending" />
    <xsl:if test="position() = 1"><xsl:value-of select="." /></xsl:if>
  </xsl:for-each>
</xsl:variable>

<main>
<table>
    <thead>
        <tr>
            <th>地址</th>
            <th class="lastmod">最后更新</th>
            <th class="changefreq">更新频率</th>
            <th class="priority">权重</th>
        </tr>
    </thead>

    <tfoot>
        <tr>
            <td colspan="1">当前总共&#160;<xsl:value-of select="count(/sm:urlset/sm:url)" />&#160;条记录</td>
            <td colspan="3" class="sitemap-lastmod">最后更新于&#160;<span><xsl:value-of select="$max" /></span></td>
        </tr>
    </tfoot>

    <tbody>
        <xsl:for-each select="sm:url">
        <!--xsl:sort select="./sm:lastmod" order="descending" /-->
        <tr>
            <td>
                <a>
                    <xsl:attribute name="href"><xsl:value-of select="sm:loc" /></xsl:attribute>
                    <xsl:value-of select="sm:loc" />
                </a>
            </td>

            <td>
                <time>
                    <xsl:value-of select="concat(substring-before(sm:lastmod, 'T'),' ',substring(sm:lastmod,12,5))" />
                </time>
            </td>

            <td>
                <xsl:choose>
                    <xsl:when test="sm:changefreq = 'never'">从不</xsl:when>
                    <xsl:when test="sm:changefreq = 'yearly'">每年</xsl:when>
                    <xsl:when test="sm:changefreq = 'monthly'">每月</xsl:when>
                    <xsl:when test="sm:changefreq = 'weekly'">每周</xsl:when>
                    <xsl:when test="sm:changefreq = 'daily'">每天</xsl:when>
                    <xsl:when test="sm:changefreq = 'hourly'">每小时</xsl:when>
                    <xsl:when test="sm:changefreq = 'always'">实时</xsl:when>
                    <xsl:otherwise><xsl:attribute name="class">error</xsl:attribute>未知的值</xsl:otherwise>
                </xsl:choose>
            </td>

            <td><xsl:value-of select="concat(sm:priority*100,'%')" /></td>
        </tr>
        </xsl:for-each>
    </tbody>
</table>
</main>
</xsl:template>

</xsl:stylesheet>
