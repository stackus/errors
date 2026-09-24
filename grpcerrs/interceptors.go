// Package grpcerrs provides gRPC interceptors that send and receive
// classified errors from github.com/stackus/errors.
//
// Install the server interceptors so every error a handler returns becomes a
// gRPC status that carries its type code, HTTP code, category, and public
// message:
//
//	server := grpc.NewServer(
//		grpc.ChainUnaryInterceptor(grpcerrs.UnaryServerInterceptor()),
//		grpc.ChainStreamInterceptor(grpcerrs.StreamServerInterceptor()),
//	)
//
// Install the client interceptors so every error a call returns is
// classified again. Pass registries holding the application kinds the client
// expects, so received errors match those kinds with errors.Is:
//
//	conn, err := grpc.NewClient(target,
//		grpc.WithChainUnaryInterceptor(grpcerrs.UnaryClientInterceptor(registry)),
//		grpc.WithChainStreamInterceptor(grpcerrs.StreamClientInterceptor(registry)),
//	)
package grpcerrs

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/stackus/errors"
)

// UnaryServerInterceptor converts the error a unary handler returns with
// [errors.SendGRPCError].
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		return resp, errors.SendGRPCError(err)
	}
}

// StreamServerInterceptor converts the error a streaming handler returns with
// [errors.SendGRPCError].
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return errors.SendGRPCError(handler(srv, ss))
	}
}

// UnaryClientInterceptor converts the error a unary call returns with
// [errors.ReceiveGRPCError], using registries to restore application kinds.
func UnaryClientInterceptor(registries ...*errors.Registry) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...), registries...)
	}
}

// StreamClientInterceptor converts the errors a streaming call returns with
// [errors.ReceiveGRPCError], using registries to restore application kinds.
// This covers the error from opening the stream and the errors from its
// Header, SendMsg, RecvMsg, and CloseSend methods. io.EOF is returned
// unchanged.
func StreamClientInterceptor(registries ...*errors.Registry) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, errors.ReceiveGRPCError(err, registries...)
		}

		return &clientStream{ClientStream: cs, registries: registries}, nil
	}
}

type clientStream struct {
	grpc.ClientStream
	registries []*errors.Registry
}

func (s *clientStream) Header() (metadata.MD, error) {
	md, err := s.ClientStream.Header()
	return md, errors.ReceiveGRPCError(err, s.registries...)
}

func (s *clientStream) CloseSend() error {
	return errors.ReceiveGRPCError(s.ClientStream.CloseSend(), s.registries...)
}

func (s *clientStream) SendMsg(m any) error {
	return errors.ReceiveGRPCError(s.ClientStream.SendMsg(m), s.registries...)
}

func (s *clientStream) RecvMsg(m any) error {
	return errors.ReceiveGRPCError(s.ClientStream.RecvMsg(m), s.registries...)
}
