package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCCode(t *testing.T) {
	t.Run("with grpc coder", func(t *testing.T) {
		assert.Equal(t, codes.NotFound, GRPCCode(ErrNotFound))
		assert.Equal(t, codes.DeadlineExceeded, GRPCCode(customTestError{codes.DeadlineExceeded, http.StatusConflict, "CUSTOM"}))
		assert.Equal(t, codes.Unavailable, GRPCCode(grpcTestError{codes.Unavailable}))
	})
	t.Run("without grpc coder", func(t *testing.T) {
		assert.Equal(t, codes.Unknown, GRPCCode(fmt.Errorf("an error")))
		assert.Equal(t, codes.Unknown, GRPCCode(httpTestError{http.StatusGatewayTimeout}))
		assert.Equal(t, codes.Unknown, GRPCCode(embedTestError{"CUSTOM"}))
		assert.Equal(t, codes.OK, GRPCCode(nil))
	})
}

func TestSendGRPCError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		err := SendGRPCError(ErrNotFound)
		s, ok := status.FromError(err)

		for _, detail := range s.Details() {
			switch d := detail.(type) {
			case *ErrorType:
				assert.Equal(t, ErrNotFound.TypeCode(), d.TypeCode)
				assert.Equal(t, ErrNotFound.GRPCCode(), codes.Code(d.GRPCCode))
				assert.Equal(t, ErrNotFound.HTTPCode(), int(d.HTTPCode))
			}
		}

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.True(t, ok)
		assert.Equal(t, 1, len(s.Details()))
		assert.Equal(t, codes.NotFound, s.Code())
		assert.Equal(t, "rpc error: code = NotFound desc = NOT_FOUND", err.Error())
		assert.Equal(t, ErrNotFound.Error(), s.Message())
	})
	t.Run("existing status", func(t *testing.T) {
		err := status.Error(codes.NotFound, "message")
		err2 := SendGRPCError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, err2)
		assert.Same(t, err, err2)
	})
	t.Run("normal error", func(t *testing.T) {
		err := SendGRPCError(fmt.Errorf("message"))
		s, ok := status.FromError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.True(t, ok)
		assert.Equal(t, codes.Unknown, s.Code())
		assert.Equal(t, "message", s.Message())

	})
	t.Run("with nil", func(t *testing.T) {
		err := SendGRPCError(nil)

		assert.Nil(t, err)
	})
	t.Run("with details", func(t *testing.T) {
		err := SendGRPCError(Wrap(ErrConflict, "conflict message"))

		s, ok := status.FromError(err)

		for _, detail := range s.Details() {
			switch d := detail.(type) {
			case *ErrorType:
				assert.Equal(t, ErrConflict.TypeCode(), d.TypeCode)
				assert.Equal(t, ErrConflict.GRPCCode(), codes.Code(d.GRPCCode))
				assert.Equal(t, ErrConflict.HTTPCode(), int(d.HTTPCode))
			}
		}

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.Equal(t, 1, len(s.Details()))
		assert.Equal(t, "rpc error: code = AlreadyExists desc = conflict message", err.Error())
		assert.Equal(t, "conflict message", s.Message())
		assert.True(t, ok)
	})
	t.Run("with custom error", func(t *testing.T) {
		testErr := customTestError{codes.NotFound, http.StatusBadRequest, "CUSTOM_ERROR"}
		err := SendGRPCError(Wrap(testErr, "custom message"))

		s, ok := status.FromError(err)

		for _, detail := range s.Details() {
			switch d := detail.(type) {
			case *ErrorType:
				assert.Equal(t, testErr.TypeCode(), d.TypeCode)
				assert.Equal(t, testErr.GRPCCode(), codes.Code(d.GRPCCode))
				assert.Equal(t, testErr.HTTPCode(), int(d.HTTPCode))
			}
		}

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.Equal(t, 1, len(s.Details()))
		assert.Equal(t, "rpc error: code = NotFound desc = custom message", err.Error())
		assert.Equal(t, "custom message", s.Message())
		assert.True(t, ok)
	})
}

func TestReceiveGRPCError(t *testing.T) {
	t.Run("is error type", func(t *testing.T) {
		err := SendGRPCError(Wrap(ErrNotFound, "message"))
		err = ReceiveGRPCError(err)

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrNotFound))
	})
	t.Run("is grpc status", func(t *testing.T) {
		err := SendGRPCError(Wrap(ErrNotFound, "message"))
		err = ReceiveGRPCError(err)
		s, ok := status.FromError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.True(t, ok)
	})
	t.Run("is custom type", func(t *testing.T) {
		testErr := customTestError{codes.Internal, http.StatusForbidden, "CONFLICT"}
		err := SendGRPCError(Wrap(testErr, "custom error"))
		err = ReceiveGRPCError(err)

		assert.NotNil(t, err)
		// pass GRPC type
		assert.True(t, stderrors.Is(err, ErrInternal))
		// pass HTTP type
		assert.True(t, stderrors.Is(err, ErrForbidden))
		// pass Embed type
		assert.True(t, stderrors.Is(err, ErrConflict))
		assert.Equal(t, "custom error", err.Error())
		assert.Equal(t, testErr.t, TypeCode(err))
		assert.Equal(t, testErr.hc, HTTPCode(err))
		assert.Equal(t, testErr.gc, GRPCCode(err))
	})
	t.Run("ok is nil", func(t *testing.T) {
		err := SendGRPCError(Wrap(ErrOK, "message"))
		err = ReceiveGRPCError(err)

		assert.Nil(t, err)
	})
	t.Run("as unknown", func(t *testing.T) {
		var err error

		err = SendGRPCError(fmt.Errorf("an error"))
		err = ReceiveGRPCError(err)
		s, ok := status.FromError(err)

		for _, detail := range s.Details() {
			switch d := detail.(type) {
			case *ErrorType:
				assert.Equal(t, ErrUnknown.TypeCode(), d.TypeCode)
				assert.Equal(t, ErrUnknown.GRPCCode(), codes.Code(d.GRPCCode))
				assert.Equal(t, ErrUnknown.HTTPCode(), int(d.HTTPCode))
			}
		}
		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrUnknown))
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.Equal(t, 1, len(s.Details()))
		assert.Equal(t, "an error", err.Error())
		assert.True(t, ok)
	})
	t.Run("sent details", func(t *testing.T) {
		err := SendGRPCError(Wrap(ErrConflict, "a conflict"))
		err = ReceiveGRPCError(err)
		s, ok := status.FromError(err)

		for _, detail := range s.Details() {
			switch d := detail.(type) {
			case *ErrorType:
				assert.Equal(t, ErrConflict.TypeCode(), d.TypeCode)
				assert.Equal(t, ErrConflict.GRPCCode(), codes.Code(d.GRPCCode))
				assert.Equal(t, ErrConflict.HTTPCode(), int(d.HTTPCode))
			}
		}

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.Equal(t, 1, len(s.Details()))
		assert.Equal(t, "a conflict", err.Error())
		assert.True(t, ok)
	})
}
