// SPDX-License-Identifier: MIT

'use strict';

const GOTO_TOP_MIN_TOP = 10;

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

createGotoTop();
