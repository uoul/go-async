# go-async

A lightweight, type-safe asynchronous execution library for Go that provides elegant abstractions for running concurrent operations with channel-based result handling.

## Features

- **Type-Safe Generics**: Leverages Go generics for compile-time type safety
- **Context-Aware**: Full support for `context.Context` for cancellation and timeout handling
- **Simple API**: Clean, intuitive interface for async operations
- **Result Pattern**: Encapsulates success and error states in a unified `Result` type
- **Single & Stream Operations**: Support for both one-shot async calls and streaming results
- **Zero Dependencies**: Built using only Go's standard library

## Installation

```bash
go get github.com/uoul/go-async
```

## Quick Start

### Single Async Operation

Execute a single asynchronous operation and receive the result through a channel:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/uoul/go-async"
)

func main() {
    ctx := context.Background()
    
    // Execute an async operation
    result := async.Do(ctx, func(ctx context.Context) (string, error) {
        time.Sleep(100 * time.Millisecond)
        return "Hello, Async!", nil
    })
    
    // Receive the result
    r := <-result
    if r.Error != nil {
        fmt.Printf("Error: %v\n", r.Error)
    } else {
        fmt.Printf("Success: %v\n", r.Value)
    }
}
```

### Streaming Results

Process multiple results asynchronously:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uoul/go-async"
)

func main() {
	v := 0
	seq := async.Stream[int](context.Background(), func(ctx context.Context) (int, error, bool) {
		v++
		time.Sleep(2 * time.Second)
		if v >= 10 {
			return v, nil, false
		}
		return v, nil, true
	})

	for result := range seq {
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			continue
		}
		fmt.Printf("Received: %v\n", result.Value)
	}
}
```

## API Reference

### Types

#### `Result[T]`
A channel type that carries async operation results:
```go
type Result[T any] chan _Result[T]
```

#### `Stream[T]`
A channel type for streaming multiple results:
```go
type Stream[T any] chan _Result[T]
```

#### `_Result[T]`
Internal result structure containing either a value or an error:
```go
type _Result[T any] struct {
    Value T
    Error error
}
```

### Functions

#### `Do[T any](ctx context.Context, action func(ctx context.Context) (T, error)) Result[T]`

Executes an action asynchronously and returns a `Result[T]` channel that receives exactly one value.

**Parameters:**
- `ctx`: Context for cancellation and timeout handling
- `action`: Function to execute asynchronously

**Returns:** A `Result[T]` channel that will receive one result and then close

#### `Stream[T any](ctx context.Context, step func(ctx context.Context) (T, error, bool)) Stream[T]`

Executes a step function repeatedly and streams results through a channel.

**Parameters:**
- `ctx`: Context for cancellation and timeout handling
- `step`: Function that returns `(value, error, shouldContinue)`

**Returns:** A `Stream[T]` channel that receives multiple results

## Design Philosophy

This library embraces Go's native concurrency primitives while providing a cleaner abstraction layer. It follows these principles:

- **Explicit over implicit**: Operations are clearly async through the API
- **Context-first**: All operations respect context cancellation
- **Type safety**: Generics ensure compile-time type checking
- **Simplicity**: Minimal API surface with maximum utility

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

---

Made with ❤️ for the Go community
