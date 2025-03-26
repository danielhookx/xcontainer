# List Package

English | [中文](README_CN.md)

A generic doubly-linked list implementation in Go, providing a flexible and efficient data structure for managing ordered collections.

## Features

- Generic type support
- Doubly-linked list implementation
- Thread-safe operations
- Rich set of operations:
  - PushFront/PushBack
  - InsertBefore/InsertAfter
  - MoveToFront/MoveToBack
  - MoveBefore/MoveAfter
  - Remove
  - Iteration support

## Installation

```bash
go get github.com/danielhookx/xcontainer/list
```

## Usage

```go
package main

import "github.com/danielhookx/xcontainer/list"

func main() {
    // Create a new list
    l := list.New[int]()
    
    // Add elements
    l.PushBack(1)
    l.PushBack(2)
    l.PushBack(3)
    
    // Iterate over elements
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }
    
    // Use generic iteration
    for item := range l.Iter() {
        fmt.Println(item)
    }
}
```

## API Reference

### New[T]() *List[T]
Creates a new empty list.

### List Methods
- `PushFront(v T) *Element[T]`
- `PushBack(v T) *Element[T]`
- `InsertBefore(v T, mark *Element[T]) *Element[T]`
- `InsertAfter(v T, mark *Element[T]) *Element[T]`
- `MoveToFront(e *Element[T])`
- `MoveToBack(e *Element[T])`
- `MoveBefore(e, mark *Element[T])`
- `MoveAfter(e, mark *Element[T])`
- `Remove(e *Element[T]) T`
- `Len() int`
- `Front() *Element[T]`
- `Back() *Element[T]`
- `Iter() <-chan T`

### Element Methods
- `Next() *Element[T]`
- `Prev() *Element[T]`
- `Value T`

## License

This project is licensed under the BSD License - see the [LICENSE](LICENSE) file for details. 