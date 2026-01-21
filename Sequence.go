package async

// Sequence is technically the same as Result, but it has other semantics
// - A Result is meant to return a single result while a Sequence is meant to return multiple
type Sequence[T any] Result[T]
