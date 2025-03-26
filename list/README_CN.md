# List 包

[English](README.md) | 中文

Go 语言实现的泛型双向链表，提供了一个灵活且高效的数据结构来管理有序集合。

## 特性

- 泛型类型支持
- 双向链表实现
- 线程安全操作
- 丰富的操作方法：
  - PushFront/PushBack（前插/后插）
  - InsertBefore/InsertAfter（指定位置前/后插入）
  - MoveToFront/MoveToBack（移动到头部/尾部）
  - MoveBefore/MoveAfter（移动到指定元素前/后）
  - Remove（删除元素）
  - 迭代器支持

## 安装

```bash
go get github.com/danielhookx/xcontainer/list
```

## 使用示例

```go
package main

import "github.com/danielhookx/xcontainer/list"

func main() {
    // 创建新链表
    l := list.New[int]()
    
    // 添加元素
    l.PushBack(1)
    l.PushBack(2)
    l.PushBack(3)
    
    // 遍历元素
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }
    
    // 使用泛型迭代器
    for item := range l.Iter() {
        fmt.Println(item)
    }
}
```

## API 参考

### New[T]() *List[T]
创建一个新的空链表。

### 链表方法
- `PushFront(v T) *Element[T]` - 在链表头部插入元素
- `PushBack(v T) *Element[T]` - 在链表尾部插入元素
- `InsertBefore(v T, mark *Element[T]) *Element[T]` - 在指定元素前插入
- `InsertAfter(v T, mark *Element[T]) *Element[T]` - 在指定元素后插入
- `MoveToFront(e *Element[T])` - 将元素移动到链表头部
- `MoveToBack(e *Element[T])` - 将元素移动到链表尾部
- `MoveBefore(e, mark *Element[T])` - 将元素移动到指定元素前
- `MoveAfter(e, mark *Element[T])` - 将元素移动到指定元素后
- `Remove(e *Element[T]) T` - 删除指定元素
- `Len() int` - 获取链表长度
- `Front() *Element[T]` - 获取链表头部元素
- `Back() *Element[T]` - 获取链表尾部元素
- `Iter() <-chan T` - 获取泛型迭代器

### 元素方法
- `Next() *Element[T]` - 获取下一个元素
- `Prev() *Element[T]` - 获取上一个元素
- `Value T` - 获取元素值

## 许可证

本项目采用 BSD 许可证 - 详见 [LICENSE](LICENSE) 文件。 