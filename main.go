package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jb843051627/strata-forge/internal/app"
)

func main() {
	database := flag.String("db", "strata-forge.db", "SQLite 文件路径")
	addr := flag.String("addr", ":8080", "HTTP 监听地址")
	smoke := flag.Bool("smoke-test", false, "执行启动探针后退出")
	flag.Parse()

	application, err := app.New(app.Config{DatabasePath: *database, Address: *addr})
	if err != nil {
		log.Printf("startup failed: %v", err)
		os.Exit(1)
	}
	defer application.Close()

	if *smoke {
		if err := application.Smoke(context.Background()); err != nil {
			log.Printf("smoke test failed: %v", err)
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := application.Run(ctx); err != nil && ctx.Err() == nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
