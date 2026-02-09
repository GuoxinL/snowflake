package snowflake

type Option struct {
	// 节点id分配器
	allocator NodeIdAllocator
	// 时间同步器
	synchronizer TimeSynchronizer
}

// OptionAppender 选项追加器
type OptionAppender interface {
	Append(op *Option)
}

// WithNodeIdAllocator 设置节点id分配器
func WithNodeIdAllocator(allocator NodeIdAllocator) func(op *Option) {
	return func(op *Option) {
		op.allocator = allocator
	}
}

// WithTimeSynchronizer 设置时间同步器
func WithTimeSynchronizer(synchronizer TimeSynchronizer) func(op *Option) {
	return func(op *Option) {
		op.synchronizer = synchronizer
	}
}
