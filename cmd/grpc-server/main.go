package main

import (
	"context"
	"flag"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/opencars/core/pkg/api/grpc"
	"github.com/opencars/core/pkg/domain/core"
	"github.com/opencars/core/pkg/domain/operation"
	"github.com/opencars/core/pkg/domain/registration"
	"github.com/opencars/core/pkg/domain/vin_decoding"

	"github.com/opencars/core/pkg/config"
	"github.com/opencars/core/pkg/logger"
)

func main() {
	cfg := flag.String("config", "config/config.yaml", "Path to the configuration file")
	port := flag.Int("port", 3000, "Port of the server")

	flag.Parse()

	conf, err := config.New(*cfg)
	if err != nil {
		logger.Fatalf("config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	r, err := registration.NewService(conf.GRPC.Registrations.Address())
	if err != nil {
		logger.Fatalf("registration service: %s", err)
	}

	o, err := operation.NewService(conf.GRPC.Operations.Address())
	if err != nil {
		logger.Fatalf("operation service: %s", err)
	}

	vd, err := vin_decoding.NewService(conf.GRPC.VinDecoding.Address())
	if err != nil {
		logger.Fatalf("vin_decoding service: %s", err)
	}

	svc, err := core.NewService(r, o, vd)
	if err != nil {
		logger.Fatalf("core service: %s", err)
	}

	addr := ":" + strconv.Itoa(*port)
	api := grpc.New(addr, svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof("Listening on %s...", addr)
	if err := api.Run(ctx); err != nil {
		logger.Fatalf("grpc: %v", err)
	}
}
