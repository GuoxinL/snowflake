# CODEBUDDY.md 本文件为 CodeBuddy 在本仓库中工作提供指导。

## 常用命令

**运行测试:** `go test`
运行包中的所有测试。

**运行特定测试:** `go test -run TestName`
按名称运行特定的测试函数。

**性能基准测试:** `go test -run=^$ -bench=.`
运行基准测试以测量生成器性能（最大约每毫秒生成 4096 个 ID）。

**生成代码覆盖率:** `go test -cover`
运行测试并生成覆盖率报告。

## 架构

这是一个基于 bwmarrin/snowflake 优化的分布式唯一 ID 生成库，解决了节点 ID 冲突和时钟回拨问题。核心架构围绕以下几个关键组件：

**Node 结构体** (`snowflake.Node`): ID 生成器实例，包含：
- `allocator`: 节点 ID 分配器接口，解决多实例节点 ID 冲突问题
- `synchronizer`: 时间同步器接口，异步上报时间防止重启后时钟回拨
- 原有的 `Generate()` 方法，在生成 ID 时调用 `synchronizer.Async(now)` 异步同步时间
- 使用互斥锁保护的内部状态（时间戳、序列号、节点 ID）

**NodeIdAllocator 接口**: 负责分配唯一的节点 ID
```go
type NodeIdAllocator interface {
    Alloc() (nodeId int64, err error)
}
```
可通过第三方存储（Etcd、Zookeeper、MySQL 等）实现节点 ID 的自动分配，支持服务名称 + IP + Port + 部署类型的键结构。

**TimeSynchronizer 接口**: 异步时间同步机制
```go
type TimeSynchronizer interface {
    Async(time int64)
}
```
采用异步方式上报当前时间戳到第三方存储，减少外部组件依赖，提高鲁棒性。

**Option 模式**: 通过 `WithNodeIdAllocator()` 和 `WithTimeSynchronizer()` 函数配置节点实例
- `NewWithOption()`: 使用选项模式创建节点，自动调用 allocator.Alloc() 获取节点 ID
- `NewNode()`: 原始构造函数，需手动指定节点 ID

**ID 类型** (`snowflake.ID`): 自定义 int64 类型，提供多种格式转换方法（Base2、Base32、Base36、Base58、Base64、字节、JSON）。ID 结构遵循 63 位格式：`[1 未使用 | 41 时间戳 | 10 节点ID | 12 序列号]`。

**配置**: 全局变量 `Epoch`、`NodeBits` 和 `StepBits` 控制 ID 格式。必须在调用 `NewNode()` 之前设置。节点和序列字段共享总共 22 位（默认 10 + 12）。

**线程安全**: 节点使用 `sync.Mutex` 实现并发安全的生成。全局包级互斥锁已弃用，将被移除。
