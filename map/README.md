# XContainer OrderedMap

English | [中文](README_CN.md)

XContainer OrderedMap is a Go implementation of an ordered map container that maintains the insertion order of key-value pairs while providing efficient lookup performance.

## Features

- Maintains insertion order of key-value pairs
- Generic support for storing any comparable key type and any value type
- Complete JSON serialization and deserialization support
- Iterator support for traversal
- Deep copy functionality
- Thread-safe (requires external synchronization for concurrent access)

## Installation

```bash
go get github.com/danielhookx/xcontainer/map
```

## Usage Examples

### Basic Operations

```go
package main

import (
    "fmt"
    "github.com/danielhookx/xcontainer/map"
)

func main() {
    // Create a new ordered map
    m := xmap.NewOrderedMap[string, int]()

    // Add key-value pairs
    m.Set("first", 1)
    m.Set("second", 2)
    m.Set("third", 3)

    // Get value
    if value, exists := m.Get("second"); exists {
        fmt.Printf("Value: %d\n", value) // Output: Value: 2
    }

    // Delete key-value pair
    m.Delete("second")

    // Get length
    fmt.Printf("Length: %d\n", m.Len()) // Output: Length: 2
}
```

### Iteration

```go
func main() {
    m := xmap.NewOrderedMap[string, string]()
    m.Set("name", "Alice")
    m.Set("age", "25")
    m.Set("city", "Beijing")

    // Iterate using the iterator
    for k, v := range m.Iter() {
        fmt.Printf("%s: %s\n", k, v)
    }
    // Output:
    // name: Alice
    // age: 25
    // city: Beijing
}
```

### JSON Serialization

```go
func main() {
    m := xmap.NewOrderedMap[string, interface{}]()
    m.Set("name", "Bob")
    m.Set("age", 30)
    m.Set("hobbies", []string{"reading", "gaming"})

    // Serialize to JSON
    jsonData, err := json.Marshal(m)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData))
    // Output: {"name":"Bob","age":30,"hobbies":["reading","gaming"]}

    // Deserialize from JSON
    newMap := xmap.NewOrderedMap[string, interface{}]()
    if err := json.Unmarshal(jsonData, newMap); err != nil {
        log.Fatal(err)
    }
}
```

### Deep Copy

```go
func main() {
    m := xmap.NewOrderedMap[string, int]()
    m.Set("a", 1)
    m.Set("b", 2)

    // Create a deep copy
    m2 := m.Copy()
    
    // Modifying the original map won't affect the copy
    m.Set("c", 3)
    fmt.Println(m2.Len()) // Output: 2
}
```

## API Reference

### Main Methods

- `NewOrderedMap[K comparable, V any]() *OrderedMap[K, V]` - Creates a new ordered map
- `Get(key K) (V, bool)` - Retrieves a value by key
- `Set(key K, value V) bool` - Sets a key-value pair
- `Delete(key K) bool` - Removes a key-value pair
- `Len() int` - Returns the number of key-value pairs
- `Copy() *OrderedMap[K, V]` - Creates a deep copy
- `ToArray() []V` - Returns all values in order
- `Iter() iter.Seq2[K, V]` - Returns an iterator

### JSON Related Methods

- `MarshalJSON() ([]byte, error)` - Serializes to JSON
- `UnmarshalJSON(data []byte) error` - Deserializes from JSON

## Notes

1. This implementation is not thread-safe and requires external synchronization for concurrent access
2. Key types must be comparable
3. All keys are converted to strings during JSON serialization
4. The iterator directly uses the underlying linked list for iteration, which means modifications during iteration will be reflected in the current iteration
