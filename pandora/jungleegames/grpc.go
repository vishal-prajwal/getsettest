package jungleegames

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	// This is imported so that we can use the grpc consul load balancer in DialGrpc
	_ "github.com/mbobakov/grpc-consul-resolver"
)

const GRPCLBPolicy = `{"loadBalancingPolicy": "round_robin"}`

// StartGrpcServer will setup a grpc server and start the grpc server as well as cleanly handle the shutdown process
//
// An example of how to use the StartGrpcServer in conjunction with jungleegames.RunAndWait
//
//	jungleegames.RunAndWait(
//		func(ctx context.Context) error {
//			err := jungleegames.StartServer("jungleegames-pandora-grpc", "localhost:9000", func (srv *grpc.Server) {
//				api.RegisterAwesomeService(srv, &service{})
//			})
//
//			return errors.Wrap(err, "failed to start grpc server")
//		}
//	)
func StartGrpcServer(grpcBindAddr string, configure func(*grpc.Server)) error {
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			sentry.CurrentHub().Recover(p)

			fmt.Println(p)

			return status.Errorf(codes.Internal, "%s", p)
		}),
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Sentry recovery of panics
			recovery.UnaryServerInterceptor(opts...),

			// newrelic monitoring
			nrgrpc.UnaryServerInterceptor(GetNewRelic()),

			// prometheus
			grpc_prometheus.UnaryServerInterceptor,

			UnaryContextSetter(),
		),

		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(opts...),

			// newrelic monitoring
			nrgrpc.StreamServerInterceptor(GetNewRelic()),

			// prometheus
			grpc_prometheus.StreamServerInterceptor,
		),

		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)

	healthServer := health.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	listener, err := net.Listen("tcp", grpcBindAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on bindAddr: %s", grpcBindAddr)
	}

	configure(server)

	shutdown.FirstFn(func() {
		server.GracefulStop()
	})

	return errors.WithStack(server.Serve(listener))
}

func GrpcCallOptions() []grpc.CallOption {
	return DefaultGrpcCallOptions()
}

func DefaultGrpcCallOptions() []grpc.CallOption {
	return []grpc.CallOption{
		grpc.CallOption(retry.WithMax(5)),
		grpc.CallOption(retry.WithPerRetryTimeout(3 * time.Second)),
		grpc.CallOption(retry.WithBackoff(retry.BackoffLinearWithJitter(3*time.Second, 2.0))),
	}
}

// DefaultDialOptions will configure for the provided service according to the default practises of jungleegames
//
// Example for connecting to grpc service in jungleegames
//
//	conn, err := grpc.Dial(consulClient.GetConnectionString("jungleegames-pandora-grpc"), jungleegames.DefaultDialOptions()...)
//	if err != nil {
//		panic(err)
//	}
func DefaultDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(GRPCLBPolicy),

		grpc.WithChainUnaryInterceptor(
			nrgrpc.UnaryClientInterceptor,
		),

		grpc.WithChainStreamInterceptor(
			nrgrpc.StreamClientInterceptor,
		),
	}
}

func GrpcDial(serviceName string) (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("%s:6000", serviceName), DefaultDialOptions()...)
}

// UnaryContextSetter map user ctx value to User object for grpc services
func UnaryContextSetter() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// if user := auth.BuildUserObject(ctx); user != nil {
		// 	ctx = context.WithValue(ctx, auth.CtxUserKey, user)
		// }

		return handler(ctx, req)
	}
}
