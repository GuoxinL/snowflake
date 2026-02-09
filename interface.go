package snowflake

// NodeIdAllocator 节点id分配器 问题1：NodeId防撞车
type NodeIdAllocator interface {
	// Alloc 分配节点ID
	Alloc() (nodeId int64, err error)
}

// TimeSynchronizer 时间同步器 问题2：防止重启时雪花ID检查大时钟回拨问题
type TimeSynchronizer interface {
	// Async 异步同步时间 该方法不报错，减少外部组件依赖增加鲁棒性
	Async(time int64)
	// Init 初始化 该方法在初始化阶段检查时间回拨问题
	Init() error
}
