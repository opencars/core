package grpc

import (
	"context"
	"net"

	"github.com/opencars/grpc/pkg/core"
	"google.golang.org/grpc"

	"github.com/opencars/core/pkg/domain"
)

// API represents gRPC API server.
type API struct {
	addr string
	s    *grpc.Server

	vehicleHandler         *vehicleHandler
	externalVehicleHandler *externalVehicleHandler
}

func New(addr string, svc domain.CoreService, external domain.ExternalService) *API {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RequestLoggingInterceptor,
		),
	}

	return &API{
		addr:                   addr,
		s:                      grpc.NewServer(opts...),
		vehicleHandler:         newVehicleHandler(svc),
		externalVehicleHandler: newExternalVehicleHandler(external),
	}
}

func (a *API) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	core.RegisterVehicleServiceServer(a.s, a.vehicleHandler)
	core.RegisterExternalVehicleServiceServer(a.s, a.externalVehicleHandler)

	errors := make(chan error)
	go func() {
		errors <- a.s.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		a.s.GracefulStop()
		return <-errors
	case err := <-errors:
		return err
	}
}
