package main

import (
	"context"
	"github.com/labstack/echo"
	"go-api-samp/util/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
 * サーバーの起動と終了シグナル受け取り後のシャットダウン
 */

func main() {
	//provider := GetProviderFactory()

	e := echo.New()
	//controller.RegisterRoute(e, providers.GetServiceProvider())
	logger := log.GetLogger()

	go func() {
		if err := e.Start(":8080"); err != http.ErrServerClosed {
			logger.Error(context.Background(), "failed to start", err.Error())
		} else {
			logger.Info(context.Background(), "shutting down")
		}
	}()

	hook := make(chan os.Signal, 1)
	signal.Notify(hook, syscall.SIGTERM, os.Interrupt)

	<-hook

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "failed to shutdown server normally: ", err.Error())
	}
}
