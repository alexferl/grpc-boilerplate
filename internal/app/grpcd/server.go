package grpcd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	"grpc-boilerplate/internal/app/grpcd/methods"
	"grpc-boilerplate/internal/pkg/logging"
)

var (
	grpcServer       *grpc.Server
	grpcHealthServer *grpc.Server
)

func init() {
	c := NewConfig()
	c.BindFlags()
	lc := logging.Config{
		LogLevel:  viper.GetString("log-level"),
		LogOutput: viper.GetString("log-output"),
		LogWriter: viper.GetString("log-writer"),
	}
	err := logging.Init(lc)
	if err != nil {
		log.Fatal().Msgf("Error initializing logger: '%v'", err)
	}
}

func Start() {
	serviceName := fmt.Sprintf("grpc.health.v1.%s", viper.GetString("app-name"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	g, ctx := errgroup.WithContext(ctx)

	// gRPC Health Server
	healthServer := health.NewServer()
	g.Go(func() error {
		addr := fmt.Sprintf("%s:%s", viper.GetString("health-bind-address"), viper.GetString("health-bind-port"))
		grpcHealthServer = grpc.NewServer()

		healthpb.RegisterHealthServer(grpcHealthServer, healthServer)

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal().Msgf("failed to listen: '%v'", err)
		}

		log.Info().Msgf("gRPC health server serving at: %s", addr)
		return grpcHealthServer.Serve(lis)
	})

	// gRPC server
	g.Go(func() error {
		addr := fmt.Sprintf("%s:%s", viper.GetString("bind-address"), viper.GetString("bind-port"))

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal().Msgf("failed to listen: '%v'", err)
		}

		service := &pb.GreeterService{
			SayHello: methods.SayHello,
		}
		grpcServer = grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
		)

		pb.RegisterGreeterService(grpcServer, service)

		log.Info().Msgf("gRPC server serving at %s", addr)

		healthServer.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), healthpb.HealthCheckResponse_SERVING)

		return grpcServer.Serve(lis)
	})

	select {
	case <-sig:
		break
	case <-ctx.Done():
		break
	}

	cancel()

	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)

	timeout := time.Duration(viper.GetInt64("graceful-timeout")) * time.Second
	_, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	if grpcHealthServer != nil {
		grpcHealthServer.GracefulStop()
	}

	err := g.Wait()
	if err != nil {
		log.Fatal().Msgf("server returned error: '%v'", err)
	}
}
