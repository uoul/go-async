package async

type _Result[T any] struct {
	Value T
	Error error
}

type Result[T any] chan _Result[T]

// This function generates a successful async result
// - Anyway this function it not necessary when using Exec() or Stream
func success[T any](val T) _Result[T] {
	return _Result[T]{
		Value: val,
		Error: nil,
	}
}

// This generates an error async result
// - Anyway this function it not necessary when using Exec() or Stream
func fail[T any](err error) _Result[T] {
	return _Result[T]{
		Value: *new(T),
		Error: err,
	}
}
