Web 后端项目：Strata Forge 负责沉积物岩芯年代学实验记录与质量复核。
# Strata Forge

这是一个单体 Go Web 服务，使用 SQLite 文件记录岩芯样品、分层、测量、复核和年代报告。

## 启动

```bash
GOTOOLCHAIN=local go run . --db strata-forge.db --addr :8080
```

服务启动后，`/healthz` 返回健康状态，根路径提供实验员操作页。
