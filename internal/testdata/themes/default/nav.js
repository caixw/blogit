// SPDX-License-Identifier: MIT

'use strict';

const GOTO_TOP_MIN_TOP = 10;

const START_HEADLINE_LEVEL = 2;

if (START_HEADLINE_LEVEL < 1 || START_HEADLINE_LEVEL > 6) {
    throw `START_HEADLINE_LEVEL 必须介于 [1,6] 之间，当前值为 ${START_HEADLINE_LEVEL}`
}

function createGotoTop() {
    const gotoTop = document.querySelector('#goto-top');

    window.addEventListener('scroll', () => {
        const scrollTop = document.body.scrollTop + document.documentElement.scrollTop;
        gotoTop.style.display = scrollTop > GOTO_TOP_MIN_TOP ? 'block' : 'none';
    }, false);

    window.addEventListener('load', () => {
        const scrollTop = document.body.scrollTop + document.documentElement.scrollTop;
        gotoTop.style.display = scrollTop > GOTO_TOP_MIN_TOP ? 'block' : 'none';
    }, false);
}

function createTOC() {
    let headlineCount = 0;

    const ids = [];
    for (let i = START_HEADLINE_LEVEL; i <= 6; i++) {
        ids.push(`h${i}`)
    }

    const elem = document.querySelector('#content');
    const headers = elem.querySelectorAll(ids.join(','));
    if (headers.length === 0) { // 没有需要处理的标签
        return;
    }

    const toc = document.querySelector('#toc');
    headers.forEach(header => {
        let id = header.getAttribute('id');
        if (!id) {
            id = `__headline--${headlineCount}`;
            headlineCount++;
            header.setAttribute('id', id);
        }

        const level = parseInt(header.tagName.charAt(1)) - START_HEADLINE_LEVEL;
        const html = `<li><a style="margin-left:${level}rem" href="#${id}">${header.innerHTML}</a></li>`
        toc.appendChild(htmlToElements(html));
    });

    document.querySelector('#toc-button').style.display = 'block';
}

// 将字符串转换成 element
function htmlToElements(html) {
    var template = document.createElement('template');
    template.innerHTML = html;
    return template.content.firstChild;
}

function nav() {
    createTOC();
    createGotoTop();
}

nav({});
