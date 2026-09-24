package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	errorCategory struct {
		httpCode int
		grpcCode codes.Code
	}

	errorClassification struct {
		kind     *Kind
		category Error

		typeCode string
		httpCode int
		grpcCode codes.Code

		publicMessage string

		// status is the gRPC status the classification was read from, if
		// any. SendGRPCError keeps its details.
		status *status.Status
	}
)

var errorCategories = map[Error]errorCategory{
	ErrCanceled:           {http.StatusRequestTimeout, codes.Canceled},
	ErrUnknown:            {http.StatusInternalServerError, codes.Unknown},
	ErrInvalidArgument:    {http.StatusBadRequest, codes.InvalidArgument},
	ErrDeadlineExceeded:   {http.StatusGatewayTimeout, codes.DeadlineExceeded},
	ErrNotFound:           {http.StatusNotFound, codes.NotFound},
	ErrAlreadyExists:      {http.StatusConflict, codes.AlreadyExists},
	ErrPermissionDenied:   {http.StatusForbidden, codes.PermissionDenied},
	ErrResourceExhausted:  {http.StatusTooManyRequests, codes.ResourceExhausted},
	ErrFailedPrecondition: {http.StatusBadRequest, codes.FailedPrecondition},
	ErrAborted:            {http.StatusConflict, codes.Aborted},
	ErrOutOfRange:         {http.StatusUnprocessableEntity, codes.OutOfRange},
	ErrUnimplemented:      {http.StatusNotImplemented, codes.Unimplemented},
	ErrInternal:           {http.StatusInternalServerError, codes.Internal},
	ErrUnavailable:        {http.StatusServiceUnavailable, codes.Unavailable},
	ErrDataLoss:           {http.StatusInternalServerError, codes.DataLoss},
	ErrUnauthenticated:    {http.StatusUnauthorized, codes.Unauthenticated},

	// HTTP-oriented categories.

	ErrBadRequest:                 {http.StatusBadRequest, codes.InvalidArgument},
	ErrUnauthorized:               {http.StatusUnauthorized, codes.Unauthenticated},
	ErrForbidden:                  {http.StatusForbidden, codes.PermissionDenied},
	ErrMethodNotAllowed:           {http.StatusMethodNotAllowed, codes.Unimplemented},
	ErrRequestTimeout:             {http.StatusRequestTimeout, codes.DeadlineExceeded},
	ErrConflict:                   {http.StatusConflict, codes.Aborted},
	ErrGone:                       {http.StatusGone, codes.NotFound},
	ErrUnsupportedMediaType:       {http.StatusUnsupportedMediaType, codes.InvalidArgument},
	ErrImATeapot:                  {http.StatusTeapot, codes.Unknown},
	ErrUnprocessableEntity:        {http.StatusUnprocessableEntity, codes.InvalidArgument},
	ErrTooManyRequests:            {http.StatusTooManyRequests, codes.ResourceExhausted},
	ErrUnavailableForLegalReasons: {http.StatusUnavailableForLegalReasons, codes.PermissionDenied},
	ErrInternalServerError:        {http.StatusInternalServerError, codes.Internal},
	ErrNotImplemented:             {http.StatusNotImplemented, codes.Unimplemented},
	ErrBadGateway:                 {http.StatusBadGateway, codes.Unavailable},
	ErrServiceUnavailable:         {http.StatusServiceUnavailable, codes.Unavailable},
	ErrGatewayTimeout:             {http.StatusGatewayTimeout, codes.DeadlineExceeded},
}

func unknownClassification() errorClassification {
	return errorClassification{
		category: ErrUnknown,
		typeCode: string(ErrUnknown),
		httpCode: http.StatusInternalServerError,
		grpcCode: codes.Unknown,
	}
}

func categoryClassification(e Error) errorClassification {
	cat, ok := errorCategories[e]
	if !ok {
		return unknownClassification()
	}

	return errorClassification{
		category: e,
		typeCode: string(e),
		httpCode: cat.httpCode,
		grpcCode: cat.grpcCode,
	}
}

func categoryByCode(code string) (Error, bool) {
	category := Error(code)

	_, ok := errorCategories[category]
	return category, ok
}
