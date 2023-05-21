package convert

// FromPtr is a safe conversion from ptr to value
// however this silents if the underlying value is nil
// thus use with caution
func FromPtr[T any](v *T) T {
	if v == nil {
		var res T
		return res
	}
	return *v
}

func ToPtr[T any](v T) *T {
	return &v
}
