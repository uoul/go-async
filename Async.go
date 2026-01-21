package async

import (
	"context"
)

// Do executes the given action asynchronously in a goroutine and returns
// a Result[T] channel that will receive exactly one value.
//
// The action is executed with the provided context, allowing for cancellation
// and timeout handling. The function returns immediately without blocking,
// and the result is sent through the returned channel once the action completes.
//
// The returned channel is closed after the result is sent, ensuring that
// consumers can safely range over it or use it in select statements.
//
// If the action returns an error, the channel receives a _Result[T] with a
// non-nil Error field and a zero-value Value field. If the action succeeds,
// the channel receives a _Result[T] with the result in the Value field and
// a nil Error field.
//
// Example:
//
//	result := Do(ctx, func(ctx context.Context) (string, error) {
//	    return fetchData(ctx)
//	})
//	r := <-result
//	if r.Error != nil {
//	    log.Printf("error: %v", r.Error)
//	} else {
//	    log.Printf("success: %v", r.Value)
//	}
func Do[T any](ctx context.Context, action func(ctx context.Context) (T, error)) Result[T] {
	r := make(Result[T])
	go func() {
		defer close(r)
		result, err := action(ctx)
		if err != nil {
			r <- fail[T](err)
		} else {
			r <- success[T](result)
		}
	}()
	return r
}

// Stream executes the given step function repeatedly in a goroutine and
// returns a Sequence[T] channel that receives multiple results.
//
// The step function is called in a loop with the provided context. Each
// iteration produces a result of type T, an error, and a boolean indicating
// whether to continue (true) or stop (false).
//
// The function returns immediately without blocking, and results are sent
// through the returned channel as they are produced. The channel is closed
// when the step function returns false or when the goroutine completes.
//
// For each iteration:
//   - If step returns an error, a _Result[T] with a non-nil Error field is sent
//   - If step succeeds, a _Result[T] with the result in the Value field is sent
//   - If the next boolean is false, the loop terminates and the channel closes
//
// This is useful for streaming operations where multiple values are produced
// over time, such as paginated API calls, database cursors, or iterative
// computations.
//
// Example:
//
//	seq := Stream(ctx, func(ctx context.Context) (int, error, bool) {
//	    value, hasMore := fetchNextBatch(ctx)
//	    return value, nil, hasMore
//	})
//	for result := range seq {
//	    if result.Error != nil {
//	        log.Printf("error: %v", result.Error)
//	        continue
//	    }
//	    log.Printf("received: %v", result.Value)
//	}
func Stream[T any](ctx context.Context, step func(ctx context.Context) (T, error, bool)) Sequence[T] {
	r := make(Sequence[T])
	go func() {
		defer close(r)
		for {
			result, err, next := step(ctx)
			if err != nil {
				r <- fail[T](err)
			} else {
				r <- success[T](result)
			}
			if !next {
				break
			}
		}
	}()
	return r
}
