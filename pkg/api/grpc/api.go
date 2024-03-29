package grpc

import (
	"context"
	"net"

	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/core/customer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/opencars/core/pkg/domain"
)

// API represents gRPC API server.
type API struct {
	addr string
	s    *grpc.Server

	vehicleHandler  *vehicleHandler
	customerHandler *customerHandler
}

func New(addr string, svc domain.CoreService, external domain.CustomerService) *API {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RequestLoggingInterceptor,
			AuthInterceptor,
		),
	}

	return &API{
		addr:            addr,
		s:               grpc.NewServer(opts...),
		vehicleHandler:  newVehicleHandler(svc),
		customerHandler: newCustomerHandler(external),
	}
}

func (a *API) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	core.RegisterVehicleServiceServer(a.s, a.vehicleHandler)
	customer.RegisterVehicleServiceServer(a.s, a.customerHandler)
	reflection.Register(a.s)

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
