package grpcerrs_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	"github.com/stackus/errors"
	"github.com/stackus/errors/grpcerrs"
)

var ErrServiceMissing = errors.NewKind(
	"SERVICE_MISSING",
	errors.ErrNotFound,
	"service not registered",
	errors.WithHTTPCode(http.StatusGone),
	errors.WithPublicMessage("That service does not exist"),
)

// healthServer fails every call with ErrServiceMissing. Watch sends one
// response before failing, so the error arrives mid-stream.
type healthServer struct {
	healthpb.UnimplementedHealthServer
}

func (healthServer) Check(_ context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return nil, ErrServiceMissing.Msgf("no service named %q", req.GetService())
}

func (healthServer) Watch(req *healthpb.HealthCheckRequest, stream grpc.ServerStreamingServer[healthpb.HealthCheckResponse]) error {
	if err := stream.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}); err != nil {
		return err
	}
	return ErrServiceMissing.Msgf("no service named %q", req.GetService())
}

// newClient starts a server with the error interceptors over an in-memory
// connection and returns a client that uses the client interceptors.
func newClient(t testing.TB, registries ...*errors.Registry) healthpb.HealthClient {
	lis := bufconn.Listen(1 << 20)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcerrs.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(grpcerrs.StreamServerInterceptor()),
	)
	healthpb.RegisterHealthServer(server, healthServer{})
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcerrs.UnaryClientInterceptor(registries...)),
		grpc.WithChainStreamInterceptor(grpcerrs.StreamClientInterceptor(registries...)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
		server.Stop()
	})

	return healthpb.NewHealthClient(conn)
}

func requireServiceMissing(t *testing.T, err error) {
	t.Helper()
	require.ErrorIs(t, err, ErrServiceMissing)
	require.ErrorIs(t, err, errors.ErrNotFound)
	require.Equal(t, "SERVICE_MISSING", errors.TypeCode(err))
	require.Equal(t, http.StatusGone, errors.HTTPCode(err))
	require.Equal(t, codes.NotFound, errors.GRPCCode(err))
	require.Equal(t, "That service does not exist", errors.PublicMessage(err))
	require.NotContains(t, err.Error(), "no service named")
}

func TestUnary(t *testing.T) {
	registry, err := errors.NewRegistry(ErrServiceMissing)
	require.NoError(t, err)

	client := newClient(t, registry)

	_, err = client.Check(context.Background(), &healthpb.HealthCheckRequest{Service: "orders"})
	requireServiceMissing(t, err)
}

func TestStream(t *testing.T) {
	registry, err := errors.NewRegistry(ErrServiceMissing)
	require.NoError(t, err)

	client := newClient(t, registry)

	stream, err := client.Watch(context.Background(), &healthpb.HealthCheckRequest{Service: "orders"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.NoError(t, err)

	_, err = stream.Recv()
	requireServiceMissing(t, err)
	require.NotEqual(t, io.EOF, err)
}

func TestWithoutRegistry(t *testing.T) {
	client := newClient(t)

	_, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{Service: "orders"})
	require.NotErrorIs(t, err, ErrServiceMissing)
	require.ErrorIs(t, err, errors.ErrNotFound)
	require.Equal(t, "SERVICE_MISSING", errors.TypeCode(err))
	require.Equal(t, http.StatusGone, errors.HTTPCode(err))
	require.Equal(t, "That service does not exist", errors.PublicMessage(err))
}

// Install the interceptors on the server and the client. Errors that handlers
// return reach callers with the same kind, codes, and public message.
func Example() {
	lis := bufconn.Listen(1 << 20)

	// Server: convert every error a handler returns.
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcerrs.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(grpcerrs.StreamServerInterceptor()),
	)
	healthpb.RegisterHealthServer(server, healthServer{})
	go func() { _ = server.Serve(lis) }()
	defer server.Stop()

	// Client: classify every error a call returns, restoring registered kinds.
	registry, _ := errors.NewRegistry(ErrServiceMissing)
	conn, _ := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcerrs.UnaryClientInterceptor(registry)),
		grpc.WithChainStreamInterceptor(grpcerrs.StreamClientInterceptor(registry)),
	)
	defer func() { _ = conn.Close() }()

	_, err := healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{Service: "orders"})

	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrServiceMissing))
	fmt.Println(errors.TypeCode(err), errors.HTTPCode(err), errors.GRPCCode(err))
	fmt.Println(errors.PublicMessage(err))
	// Output:
	// rpc error: code = NotFound desc = That service does not exist
	// true
	// SERVICE_MISSING 410 NotFound
	// That service does not exist
}
