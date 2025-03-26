# xcontainer

## Description

xcontainer is a comprehensive collection of fundamental data structures implemented in Go. This library provides efficient and type-safe implementations of various data structures using Go's generics feature. It serves as a reliable foundation for building complex data processing applications.

## Features

- Generic-based implementation
- Type-safe data structure operations
- Clean and intuitive API design
- Comprehensive unit test coverage
- Non-thread-safe implementations (thread safety should be handled by the caller if needed)

## Data Structures

- **graph**: Graph data structure supporting directed and undirected graphs
- **heap**: Heap implementation with support for max and min heaps
- **[list](list/README.md)**: Linked list implementations including singly and doubly linked lists
- **[ordered map](map/README.md)**: Ordered map implementation maintaining insertion order of key-value pairs
- **queue**: Queue implementations including standard and priority queues
- **[set](set/README.md)**: Set implementation with basic set operations
- **stack**: Stack implementation with standard stack operations
- **tree**: Tree implementations including binary trees and binary search trees

## Requirements

- Go 1.23 or higher
- All implementations are non-thread-safe
- Generic-based implementation for type safety
- Review specific data structure documentation before use

## License

xcontainer is under the MIT license. See the [LICENSE](LICENSE) file for details.