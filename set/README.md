English | [中文](README_CN.md)

# Set Collection

Set is a generic implementation of the set data structure that supports any comparable type.

## Key Features

- Generic support for any comparable type
- Capacity limit support
- Rich set operations
- Iterator support
- Not thread-safe

## Basic Operations

```go
// Create a set
s := set.CreateSet[int]()

// Add element
s.Add(1)

// Check if element exists
exists := s.IsElementOf(1)

// Get set size
size := s.Size()

// Remove element
s.Remove(1)

// Clear set
s.Clear()
```

## Advanced Operations

```go
// Create a set with capacity
s := set.CreateSetWithCapacity[int](10)

// Create set from multiple elements
s := set.BuildSet(1, 2, 3)

// Create set from iterator
s := set.Collect(someIterator)

// Filter set
filtered := s.Filter(func(x int) bool { return x > 0 })

// Map set
mapped := s.Map(func(x int) int { return x * 2 })

// Check if set is empty
isEmpty := set.IsEmptySet(s)

// Check if two sets are equal
equal := s1.Equal(s2)
```

## Iteration

```go
// Iterate using iterator
for x := range s.Iter() {
    fmt.Println(x)
}

// Get all elements as list
elements := s.Enumerate()
```

## Notes

- Sets are not thread-safe, locking is required in concurrent environments
- When capacity limit is set, exceeding it will return an error
- Element order in sets is not guaranteed 