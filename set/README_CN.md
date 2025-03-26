[English](README.md) | 中文

# Set 集合

Set 是一个基于泛型实现的集合数据结构，支持任意可比较类型。

## 主要特性

- 支持泛型，可以存储任意可比较类型
- 支持容量限制
- 提供丰富的集合操作
- 支持迭代器
- 线程不安全

## 基本操作

```go
// 创建集合
s := set.CreateSet[int]()

// 添加元素
s.Add(1)

// 检查元素是否存在
exists := s.IsElementOf(1)

// 获取集合大小
size := s.Size()

// 删除元素
s.Remove(1)

// 清空集合
s.Clear()
```

## 高级操作

```go
// 创建带容量的集合
s := set.CreateSetWithCapacity[int](10)

// 从多个元素创建集合
s := set.BuildSet(1, 2, 3)

// 从迭代器创建集合
s := set.Collect(someIterator)

// 过滤集合
filtered := s.Filter(func(x int) bool { return x > 0 })

// 映射集合
mapped := s.Map(func(x int) int { return x * 2 })

// 检查集合是否为空
isEmpty := set.IsEmptySet(s)

// 检查两个集合是否相等
equal := s1.Equal(s2)
```

## 迭代

```go
// 使用迭代器遍历集合
for x := range s.Iter() {
    fmt.Println(x)
}

// 获取所有元素列表
elements := s.Enumerate()
```

## 注意事项

- 集合不是线程安全的，需要在并发环境下自行加锁
- 当设置了容量限制时，超出容量会返回错误
- 集合中的元素顺序是不确定的
