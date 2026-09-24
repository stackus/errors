package errors_test

import (
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestWrappedErrorCodersAndIdentity(t *testing.T) {
	kind := errors.NewKind("MISSING_RECORD", errors.ErrNotFound, "record missing")
	otherKind := errors.NewKind("OTHER_MISSING_RECORD", errors.ErrNotFound, "other record missing")
	type testCase struct {
		err        error
		typeCode   string
		httpCode   int
		grpcCode   codes.Code
		matches    []error
		nonMatches []error
	}
	tests := map[string]testCase{
		"category wrapper": {
			err:      errors.ErrNotFound.Msg("missing"),
			typeCode: "NOT_FOUND", httpCode: http.StatusNotFound, grpcCode: codes.NotFound,
			matches: []error{errors.ErrNotFound}, nonMatches: []error{errors.ErrForbidden, kind},
		},
		"kind wrapper": {
			err:      kind.Msg("missing"),
			typeCode: "MISSING_RECORD", httpCode: http.StatusNotFound, grpcCode: codes.NotFound,
			matches: []error{kind, errors.ErrNotFound}, nonMatches: []error{otherKind, errors.ErrForbidden},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, "missing", tt.err.Error())
			require.NoError(t, errors.Unwrap(tt.err))
			require.Equal(t, tt.typeCode, tt.err.(errors.TypeCoder).TypeCode())
			require.Equal(t, tt.httpCode, tt.err.(errors.HTTPCoder).HTTPCode())
			require.Equal(t, tt.grpcCode, tt.err.(errors.GRPCCoder).GRPCCode())
			for _, target := range tt.matches {
				require.ErrorIs(t, tt.err, target)
			}
			for _, target := range tt.nonMatches {
				require.NotErrorIs(t, tt.err, target)
			}
		})
	}
}
