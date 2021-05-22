package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "message")

		assert.Nil(t, err)
	})
	t.Run("suffixing", func(t *testing.T) {
		err := fmt.Errorf("suffix")
		err = Wrap(err, "prefix")

		assert.NotNil(t, err)
		assert.Equal(t, "prefix: suffix", err.Error())
	})

	t.Run("no suffixing", func(t *testing.T) {
		err := Wrap(ErrNotFound, "message")

		assert.NotNil(t, err)
		assert.Equal(t, "message", err.Error())
	})

	t.Run("multiple", func(t *testing.T) {
		err := Wrap(ErrNotFound, "suffix")
		err = Wrap(err, "prefix")

		assert.NotNil(t, err)
		assert.Equal(t, "prefix: suffix", err.Error())
	})
}

func TestWrapf(t *testing.T) {
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

	t.Run("multiple", func(t *testing.T) {
		err := Wrapf(ErrNotFound, "suffix %s", "2")
		err = Wrapf(err, "prefix %s", "2")

		assert.NotNil(t, err)
		assert.Equal(t, "prefix 2: suffix 2", err.Error())
	})
}

func ExampleMessage() {
	err := Wrap(ErrNotFound, "message")

	fmt.Println(Message(err))
	// Output: NOT_FOUND: message
}

func TestMessage(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		err := Wrap(ErrNotFound, "message")

		assert.NotNil(t, err)
		assert.Equal(t, "NOT_FOUND: message", Message(err))
	})
	t.Run("without error", func(t *testing.T) {
		err := fmt.Errorf("suffix")
		err = Wrap(err, "prefix")

		assert.NotNil(t, err)
		assert.Equal(t, "prefix: suffix", Message(err))
	})
	t.Run("with multiple errors", func(t *testing.T) {
		err := Wrap(ErrNotFound, "suffix")
		err = Wrap(err, "prefix")

		assert.NotNil(t, err)
		assert.Equal(t, "NOT_FOUND: prefix: suffix", Message(err))
	})
}
