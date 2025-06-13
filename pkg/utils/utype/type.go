package utype

// ValueToPtr returns a pointer to the value of type T.
func ValueToPtr[T any](value T) *T {
	return &value
}

// PtrToValue returns the value of type T from a T pointer.
func PtrToValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T

		return zero
	}

	return *ptr
}
