package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func ExampleWrap() {
	err := Wrap(ErrNotFound, "message")
	fmt.Println(err)
	// Output: message
}

func ExampleWrap_multiple() {
	err := Wrap(ErrNotFound, "original message")
	err = Wrap(err, "prefixed message")
	fmt.Println(err)
	// Output: prefixed message: original message
}

func TestEmbedType(t *testing.T) {
	t.Run("with embed typer", func(t *testing.T) {
		assert.Equal(t, "NOT_FOUND", TypeCode(ErrNotFound))
		assert.Equal(t, "CUSTOM", TypeCode(customTestError{codes.Internal, http.StatusConflict, "CUSTOM"}))
		assert.Equal(t, "VERY CUSTOM", TypeCode(embedTestError{"VERY CUSTOM"}))
	})
	t.Run("without embed typer", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN", TypeCode(fmt.Errorf("an error")))
		assert.Equal(t, "UNKNOWN", TypeCode(grpcTestError{codes.Canceled}))
		assert.Equal(t, "UNKNOWN", TypeCode(httpTestError{http.StatusBadGateway}))
		assert.Equal(t, "", TypeCode(nil))

	})
}

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "message")

		assert.Nil(t, err)
	})
	t.Run("embed error type", func(t *testing.T) {
		err := Wrap(ErrNotFound, "error message")

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrNotFound))
		assert.Equal(t, "error message", err.Error())
	})
	t.Run("normal error", func(t *testing.T) {
		err := Wrap(fmt.Errorf("normal error"), "message")

		assert.NotNil(t, err)
		assert.Equal(t, "message: normal error", err.Error())
	})
	t.Run("no message modification", func(t *testing.T) {
		err := Wrap(ErrNotFound, "message")

		assert.NotNil(t, err)
		assert.Equal(t, "message", err.Error())
	})
	t.Run("message modification", func(t *testing.T) {
		err := Wrap(ErrConflict, "first message")
		err = Wrap(err, "second message")

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrConflict))
		assert.Equal(t, "second message: first message", err.Error())
	})
	t.Run("multiple errors", func(t *testing.T) {
		err := Wrap(ErrNotFound, "suffix")
		err = Wrap(err, "prefix")
		err = Wrap(err, "another")

		assert.NotNil(t, err)
		assert.True(t, stderrors.Is(err, ErrNotFound))
		assert.Equal(t, "another: prefix: suffix", err.Error())
	})
	t.Run("embed custom error", func(t *testing.T) {
		testErr := customTestError{codes.NotFound, http.StatusBadRequest, "CUSTOM_ERROR"}
		err := Wrap(testErr, "custom message")

		assert.NotNil(t, err)
		assert.Equal(t, "custom message", err.Error())
	})
	t.Run("embed type error", func(t *testing.T) {
		testErr := embedTestError{"EMBED_TYPE"}
		err := Wrap(testErr, "embed message")

		assert.NotNil(t, err)
		assert.Equal(t, "embed message", err.Error())
		assert.True(t, stderrors.Is(err, testErr))
	})
	t.Run("embed http error", func(t *testing.T) {
		testErr := httpTestError{http.StatusConflict}
		err := Wrap(testErr, "http message")

		assert.NotNil(t, err)
		assert.Equal(t, "http message: HTTP(409)", err.Error())
		assert.True(t, stderrors.Is(err, testErr))
	})
	t.Run("embed grpc error", func(t *testing.T) {
		testErr := grpcTestError{codes.AlreadyExists}
		err := Wrap(testErr, "grpc message")

		assert.NotNil(t, err)
		assert.Equal(t, "grpc message: GRPC(6)", err.Error())
		assert.True(t, stderrors.Is(err, testErr))
	})
}

func TestWrapf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrapf(nil, "message")

		assert.Nil(t, err)
	})
	t.Run("suffixing", func(t *testing.T) {
		err := fmt.Errorf("suffix")
		err = Wrapf(err, "prefix %s", "2")

		assert.NotNil(t, err)
		assert.Equal(t, "prefix 2: suffix", err.Error())
	})
	t.Run("no suffixing", func(t *testing.T) {
		err := Wrapf(ErrNotFound, "message %s", "2")

		assert.NotNil(t, err)
		assert.Equal(t, "message 2", err.Error())
	})
	t.Run("multiple errors", func(t *testing.T) {
		err := Wrapf(ErrNotFound, "suffix %s", "1")
		err = Wrapf(err, "prefix %s", "2")
		err = Wrapf(err, "another %s", "3")

		assert.NotNil(t, err)
		assert.Equal(t, "another 3: prefix 2: suffix 1", err.Error())
	})
}
