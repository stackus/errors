package errors_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestClassificationOfContextAndJoinedErrors(t *testing.T) {
	type testCase struct {
		err      error
		typeCode string
		httpCode int
		grpcCode codes.Code
	}
	tests := map[string]testCase{
		"canceled context": {
			context.Canceled, "CANCELED", http.StatusRequestTimeout, codes.Canceled,
		},
		"deadline exceeded": {
			context.DeadlineExceeded, "DEADLINE_EXCEEDED", http.StatusGatewayTimeout, codes.DeadlineExceeded,
		},
		"wrapped context error": {
			errors.Wrap(context.Canceled, "operation canceled"), "CANCELED", http.StatusRequestTimeout, codes.Canceled,
		},
		"joined identical categories": {
			errors.Join(errors.ErrNotFound, errors.ErrNotFound.Msg("missing record")), "NOT_FOUND", http.StatusNotFound, codes.NotFound,
		},
		"joined classified and ordinary errors": {
			errors.Join(errors.ErrNotFound, errors.New("database unavailable")), "UNKNOWN", http.StatusInternalServerError, codes.Unknown,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.typeCode, errors.TypeCode(tt.err))
			require.Equal(t, tt.httpCode, errors.HTTPCode(tt.err))
			require.Equal(t, tt.grpcCode, errors.GRPCCode(tt.err))
		})
	}
}
