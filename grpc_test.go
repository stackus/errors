package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSendGRPCError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		err := SendGRPCError(ErrNotFound)
		s, ok := status.FromError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, s.Code())
		assert.Equal(t, ErrNotFound.Error(), s.Message())
	})
	t.Run("existing status", func(t *testing.T) {
		err := status.Error(codes.NotFound, "message")
		err2 := SendGRPCError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, err2)
		assert.Same(t, err, err2)
	})
	t.Run("as unknown", func(t *testing.T) {
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
}

func TestReceiveGRPCError(t *testing.T) {
	t.Run("is error", func(t *testing.T) {
		err := status.Error(codes.NotFound, "message")
		err = ReceiveGRPCError(err)

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrNotFound))
	})
	t.Run("is status", func(t *testing.T) {
		err := status.Error(codes.NotFound, "message")
		err = ReceiveGRPCError(err)
		s, ok := status.FromError(err)

		assert.NotNil(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &status.Status{}, s)
		assert.True(t, ok)
	})
	t.Run("ok is nil", func(t *testing.T) {
		err := status.Error(codes.OK, "message")

		err = ReceiveGRPCError(err)

		assert.Nil(t, err)
	})
	t.Run("as unknown", func(t *testing.T) {
		var err error

		err = ErrNotFound
		err = ReceiveGRPCError(err)

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrUnknown))
		assert.False(t, stderrors.Is(err, ErrNotFound))
	})
}
