本项目实现了 GroupCache 的大部分功能，比如 LRU 缓存淘汰策略、单机并发缓存、HTTP 服务端、一致性哈希(Hash)、分布式节点、防止缓存击穿和使用 Protobuf 通信。

整体上支持特性有：1.单机缓存和基于 HTTP 的分布式缓存；2.最近最少访问(Least Recently Used, LRU) 缓存策略；3.使用 Go 锁机制防止缓存击穿；4.使用一致性哈希(Hash)选择节点，实现负载均衡；5.使用 protobuf 优化节点间二进制通信。
