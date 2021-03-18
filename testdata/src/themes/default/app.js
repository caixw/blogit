// SPDX-License-Identifier: MIT

"use strict";

/**
 * 以下是一些公用常量的定义，方便修改配置。
 */
const gototopTimes = 20
const gototopDistance = 30

/**
 * 程序的入口代码。
 * DOM 加载完成之后，才执行操作，相当于 jQuery.ready
 */
window.document.addEventListener('DOMContentLoaded', ()=>{
    initTop()
})

/**
 * 初始化所有与 goto top 相关的内容
 */
function initTop () {
    const top = window.document.getElementById('top')
    const gototopDisplay = top.style.display;

    const scrollTop = function() {
        return document.body.scrollTop + document.documentElement.scrollTop
    }

    // 根据与页面顶部的距离，控制是否显示 top 按钮。
    const showTopButton = function() {
        const dsp = scrollTop() > gototopDistance ? gototopDisplay : 'none'
        if (dsp !== top.style.display) {
            top.style.display = dsp
        }
    }

    showTopButton()
    window.addEventListener('scroll', ()=>{showTopButton()})

    // 滚动到顶部
    top.addEventListener('click', (event)=>{
        let height = scrollTop()
        const offset = height / gototopTimes

        const tick = window.setInterval(function(){
            height -= offset
            window.scrollTo(0, height)

            if (height <= 0){
                window.clearInterval(tick)
            }
        }, 10)

        event.preventDefault(); // 防止继续执行 A 标签的 href
    })
}
