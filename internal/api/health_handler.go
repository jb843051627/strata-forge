package api

import (
	"context"
	"net/http"
)

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := a.service.ListSamples(context.Background(), ""); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><title>Strata Forge</title><style>body{font-family:system-ui;max-width:920px;margin:2rem auto;padding:0 1rem;color:#25323c}header{border-bottom:2px solid #c77d42}table{border-collapse:collapse;width:100%;margin-top:1rem}td,th{border:1px solid #ccd3d8;padding:.5rem;text-align:left}</style></head><body><header><h1>Strata Forge</h1><p>岩芯年代学实验记录</p></header><main><p>样品接收与测量状态</p><div id="list">加载中...</div></main><script>fetch('/api/v1/samples').then(r=>r.json()).then(rows=>document.querySelector('#list').innerHTML='<table><tr><th>编号</th><th>地点</th><th>状态</th></tr>'+rows.map(x=>'<tr><td>'+x.code+'</td><td>'+x.site+'</td><td>'+x.status+'</td></tr>').join('')+'</table>')</script></body></html>`))
}
