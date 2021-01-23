package main

import (
	"context"
	"github.com/labstack/echo"
	"go-api-samp/controller"
	"go-api-samp/util/config"
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
	provider := GetProviderFactory()

	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	e := echo.New()
	controller.RegisterRoute(e, provider.GetServiceProvider())

	logger := log.GetLogger()

	go func() {
		if err := e.Start(config.Server.Addr); err != http.ErrServerClosed {
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
