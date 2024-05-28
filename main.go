package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	echoPrometheus "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"go.uber.org/zap"

	"github.com/fernandonogueira/statuspage-exporter/pkg/config"
	"github.com/fernandonogueira/statuspage-exporter/pkg/prober"
)

const (
	shutdownTimeout = 5 * time.Second
)

func handleHealthz(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func startHTTP(ctx context.Context, wg *sync.WaitGroup, config *config.ExporterConfig) {
	wg.Add(1)
	defer wg.Done()

	srv := echo.New()
	echoPrometheus := echoPrometheus.NewPrometheus("statuspage_exporter", nil)
	echoPrometheus.Use(srv)

	srv.GET("/probe", prober.Handler(config))
	srv.GET("/healthz", handleHealthz)

	httpPort := config.HTTPPort
	httpAddr := fmt.Sprintf(":%d", httpPort)

	// Start your http server for prometheus.
	go func() {
		if err := srv.Start(httpAddr); err != nil {
			config.Log.Panic("Unable to start a http server.", zap.Error(err))
		}
	}()

	config.Log.Info("Http server listening on", zap.Int("port", httpPort))

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		config.Log.Panic("Http server Shutdown Failed", zap.Error(err))
	}

	config.Log.Info("Http server stopped")
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.Parse()

	cfg, err := config.InitConfig(configFile)

	cfg.Log.Info("Loaded config",
		zap.Int("http_port", cfg.HTTPPort),
		zap.Duration("client_timeout", cfg.ClientTimeout),
		zap.Int("retry_count", cfg.RetryCount))

	cfg.Log.Info("Starting statuspage_exporter...")

	if err != nil {
		cfg.Log.Fatal("Unable to initialize config", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	prometheus.MustRegister(version.NewCollector("statuspage_exporter"))

	go startHTTP(ctx, wg, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cfg.Log.Info("Received shutdown signal. Waiting for workers to terminate...")
	cancel()

	wg.Wait()
}
