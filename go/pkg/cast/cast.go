package cast

// Ptr returns pointer of v.
func Ptr[T any](v T) *T {
	return &v
}

// Value returns value of pointer.
// Returns a zero value if pointer is nil.
func Value[T any](p *T) T { // nolint:ireturn
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
