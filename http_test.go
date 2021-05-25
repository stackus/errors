package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestHTTPCode(t *testing.T) {
	t.Run("with http coder", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, HTTPCode(ErrNotFound))
		assert.Equal(t, http.StatusConflict, HTTPCode(customTestError{codes.Internal, http.StatusConflict, "CUSTOM"}))
		assert.Equal(t, http.StatusBadGateway, HTTPCode(httpTestError{http.StatusBadGateway}))
	})
	t.Run("without http coder", func(t *testing.T) {
		assert.Equal(t, http.StatusNotExtended, HTTPCode(fmt.Errorf("an error")))
		assert.Equal(t, http.StatusNotExtended, HTTPCode(grpcTestError{codes.Canceled}))
		assert.Equal(t, http.StatusNotExtended, HTTPCode(embedTestError{"CUSTOM"}))
		assert.Equal(t, http.StatusOK, HTTPCode(nil))
	})
}
