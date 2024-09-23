package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/damonto/estkme-cloud/internal/cloud"
	"github.com/damonto/estkme-cloud/internal/config"
)

var Version string

func init() {
	flag.StringVar(&config.C.ListenAddress, "listen-address", ":1888", "eSTK.me cloud enhance server listen address")
	flag.StringVar(&config.C.Prompt, "prompt", "", "prompt message to show on the server (max: 100 characters)")
	flag.BoolVar(&config.C.Verbose, "verbose", false, "verbose mode")
	flag.Parse()
}

func initApp() {
	config.C.LoadEnv()
	if config.C.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Warn("verbose mode is enabled, this will print out sensitive information")
	}
	if err := config.C.IsValid(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
}

func main() {
	slog.Info("eSTK.me cloud enhance server", "version", Version)
	initApp()

	manager := cloud.NewManager()
	server := cloud.NewServer(manager)

	go func() {
		if err := server.Listen(config.C.ListenAddress); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")
	if err := server.Shutdown(); err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}
}
