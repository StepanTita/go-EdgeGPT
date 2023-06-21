package iter

func Map[I, O any](values []I, mapper func(I) O) []O {
	out := make([]O, len(values))
	for i, v := range values {
		out[i] = mapper(v)
	}
	return out
}
