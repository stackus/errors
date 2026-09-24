package errors

import (
	"fmt"
)

type Registry struct {
	kinds map[string]*Kind
}

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

func (r *Registry) Lookup(code string) (*Kind, bool) {
	if r == nil {
		return nil, false
	}

	kind, ok := r.kinds[code]

	return kind, ok
}
