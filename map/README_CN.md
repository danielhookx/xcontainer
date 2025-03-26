# XContainer OrderedMap

[English](README.md) | 中文

XContainer OrderedMap 是一个 Go 语言实现的有序 Map 容器，它保持了键值对的插入顺序，同时提供了高效的查找性能。

## 特性

- 保持键值对的插入顺序
- 支持泛型，可以存储任意可比较类型的键和任意类型的值
- 提供完整的 JSON 序列化和反序列化支持
- 支持迭代器遍历
- 提供深拷贝功能

## 安装

```bash
go get github.com/danielhookx/xcontainer/map
```

## 使用示例

### 基本操作

```go
package main

import (
    "fmt"
    "github.com/danielhookx/xcontainer/map"
)

func main() {
    // 创建一个新的有序 Map
    m := xmap.NewOrderedMap[string, int]()

    // 添加键值对
    m.Set("first", 1)
    m.Set("second", 2)
    m.Set("third", 3)

    // 获取值
    if value, exists := m.Get("second"); exists {
        fmt.Printf("Value: %d\n", value) // 输出: Value: 2
    }

    // 删除键值对
    m.Delete("second")

    // 获取长度
    fmt.Printf("Length: %d\n", m.Len()) // 输出: Length: 2
}
```

### 迭代遍历

```go
func main() {
    m := xmap.NewOrderedMap[string, string]()
    m.Set("name", "Alice")
    m.Set("age", "25")
    m.Set("city", "Beijing")

    // 使用迭代器遍历
    for k, v := range m.Iter() {
        fmt.Printf("%s: %s\n", k, v)
    }
    // 输出:
    // name: Alice
    // age: 25
    // city: Beijing
}
```

### JSON 序列化

```go
func main() {
    m := xmap.NewOrderedMap[string, interface{}]()
    m.Set("name", "Bob")
    m.Set("age", 30)
    m.Set("hobbies", []string{"reading", "gaming"})

    // 序列化为 JSON
    jsonData, err := json.Marshal(m)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData))
    // 输出: {"name":"Bob","age":30,"hobbies":["reading","gaming"]}

    // 从 JSON 反序列化
    newMap := xmap.NewOrderedMap[string, interface{}]()
    if err := json.Unmarshal(jsonData, newMap); err != nil {
        log.Fatal(err)
    }
}
```

### 深拷贝

```go
func main() {
    m := xmap.NewOrderedMap[string, int]()
    m.Set("a", 1)
    m.Set("b", 2)

    // 创建深拷贝
    m2 := m.Copy()
    
    // 修改原 Map 不会影响拷贝
    m.Set("c", 3)
    fmt.Println(m2.Len()) // 输出: 2
}
```

## API 参考

### 主要方法

- `NewOrderedMap[K comparable, V any]() *OrderedMap[K, V]` - 创建新的有序 Map
- `Get(key K) (V, bool)` - 获取键对应的值
- `Set(key K, value V) bool` - 设置键值对
- `Delete(key K) bool` - 删除键值对
- `Len() int` - 获取 Map 长度
- `Copy() *OrderedMap[K, V]` - 创建深拷贝
- `ToArray() []V` - 获取所有值（保持顺序）
- `Iter() iter.Seq2[K, V]` - 获取迭代器

### JSON 相关方法

- `MarshalJSON() ([]byte, error)` - 序列化为 JSON
- `UnmarshalJSON(data []byte) error` - 从 JSON 反序列化

## 注意事项

1. 该实现不是线程安全的，在并发环境下需要外部同步
2. 键类型必须是可比较的（comparable）
3. 在 JSON 序列化时，所有键都会被转换为字符串
4. 迭代器直接使用底层链表进行迭代，这意味着在迭代过程中的修改会反映在当前迭代中
