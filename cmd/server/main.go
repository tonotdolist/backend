package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tonotdolist/cmd/server/wire"
)

func main() {
	app := wire.InitializeApp()
	app.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Stop(ctx)
}
