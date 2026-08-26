package report

import "strings"

const pageTemplate = `<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8"><title>Strata Forge</title>
<style>body{font-family:system-ui,sans-serif;max-width:960px;margin:2rem auto;padding:0 1rem;color:#23313d}header{border-bottom:2px solid #c77d42;margin-bottom:1rem}table{border-collapse:collapse;width:100%}td,th{border:1px solid #cad2d8;padding:.5rem;text-align:left}.muted{color:#6c7a83}</style></head>
<body><header><h1>Strata Forge</h1><p class="muted">岩芯年代学实验记录</p></header><main id="app">正在加载...</main>
<script>fetch('/api/v1/samples').then(r=>r.json()).then(rows=>{const main=document.querySelector('#app');main.innerHTML='<h2>最近样品</h2><table><tr><th>编号</th><th>地点</th><th>状态</th></tr>'+rows.map(function(x){return '<tr><td>'+x.code+'</td><td>'+x.site+'</td><td>'+x.status+'</td></tr>'}).join('')+'</table>'}).catch(e=>document.querySelector('#app').textContent=e)</script></body></html>`

func Page() string {
	return strings.TrimSpace(pageTemplate)
}
