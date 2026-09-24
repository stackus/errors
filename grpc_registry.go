package errors

import (
	"fmt"
)

// Registry maps type codes to the [Kind] values a gRPC client expects to
// receive. [ReceiveGRPCError] uses it so that an error sent as a Kind matches
// the client's copy of that Kind with [Is].
//
// A client that shares error definitions with its servers, for example
// through a common package, registers those kinds once at startup. A
// received kind is matched by type code, gRPC code, and category. The local
// Kind's definition supplies the HTTP code.
//
// A registry is not needed to match built-in categories, or to keep an
// unregistered kind's type code, category, and HTTP code. A Registry is safe
// for concurrent use.
type Registry struct {
	kinds map[string]*Kind
}

// NewRegistry returns a registry that contains kinds. It returns an error for
// a nil Kind or when two kinds share a type code.
func NewRegistry(kinds ...*Kind) (*Registry, error) {
	r := &Registry{
		kinds: make(map[string]*Kind, len(kinds)),
	}

	for _, kind := range kinds {
		if kind == nil {
			return nil, fmt.Errorf("errors: registry contains nil kind")
		}

		if _, exists := r.kinds[kind.typeCode]; exists {
			return nil, fmt.Errorf("errors: registry contains duplicate kind with type code %q", kind.typeCode)
		}

		r.kinds[kind.typeCode] = kind
	}

	return r, nil
}

// Lookup returns the Kind registered for code. A nil Registry is safe to use
// and returns no match.
func (r *Registry) Lookup(code string) (*Kind, bool) {
	if r == nil {
		return nil, false
	}

	kind, ok := r.kinds[code]

	return kind, ok
}
