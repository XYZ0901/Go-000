# 学习笔记
## 作业
**在work目录下哦😉**
## 分布式缓存 & 分布式事务
- 核心中间件: 缓存、存储、队列

### 缓存选型
#### memcache
- memcache提供简单的kv cache存储, value大小不超过1mb
- 使用memcache作为大文本或者简单的kv结构
- memcache使用了slab方式做内存管理, 存在一定浪费, 如果大量接近的item, 
  需要调整slab的ratio参数, 防止某些slab热点导致内存足够的情况下引发LRU.
- 大部分情况下, 简单kv推荐使用Memcache, 吞吐和响应都足够好

#### redis
- redis有丰富的数据类型, 支持增量的方式修改部分数据, 排行榜 集合 数组等
- redis没有使用内存池, 存在一定的内存碎片

#### redis VS memcache
- redis单线程, memcache多线程, 压测结果QPS差异不大, 但吞吐会有很大差别, 
  如大数据的value返回时, redis qps会厉害抖动下降, 因为单线程工作, 其它查询进不来
- 建议纯kv走memcache, 在系统中建议使用memcache + redis双缓存设计

#### Proxy
- twemproxy
  - 单进程单线程模型, 有io瓶颈
  - 二次开发成本难度高, 难以与公司运维平台集成
  - 不支持自动伸缩
  - 运维不友好
- 不推荐使用

#### Hash
- 均衡问题

#### 一致性hash
- 有界负责一致性hash
- 虚拟节点 murmurhash3

#### Slot
- CRC16(key) & 16383
- redis cluster

### 缓存模式
#### 数据一致性
- 由于Storage和Cache同步更新容易出现数据不一致
  - 同步操作DB
  - 同步操作Cache
  - 利用Job消费信息, 重新补偿一次缓存操作
  - 保证实效性和一致性
#### 多级缓存
- 保证多级缓存一致性
  - 优先清理上游再清理下游
  - 下游的缓存expire要大于上游, 里面穿透回源
#### 热点缓存
- 小表广播, 从remoteCache提升为LocalCache
- 主动监控防御预热, 比如直播房间页高在线情况下直接外挂服务主动防御
- 基础框架支持热点发现, 自动短时的short live cache
- 多Cluster支持
  - 多key设计: 使用多副本, 减小节点热点的问题
  - 空间换时间
  - 当业务频繁更新时, cache频繁过期, 导致命中率低: stale sets
#### 穿透缓存
- single fly
- 分布式锁
- 队列
- lease

### 缓存技巧
#### Incast Congestion
- 如果网络中存在大量包会出现延迟问题
#### 小技巧
- key尽可能小, 可以int绝不string
- 拆分key
- 空缓存设置, 对于部分数据由于数据库始终为空, 此时应该设置空缓存, 避免缓存miss穿透
- 空缓存保护策略
- 读失败以后的写缓存策略(降级后一般读失败不触发回写缓存)
- 序列化使用protobuf, 尽可能减少size
- 工具化浇水代码
- flag的使用 - 标识compress, encoding, large value
- memcache支持gets 尽可能减少pipeline, 减少网络往返
- 使用二进制协议, 支持pipeline delete, UDP读取, TCP更新
#### redis 小技巧
- 增量更新一致性
- BITSET: 存储每日登陆用户, 单个标记位置(boolean), 为避免单个BITSET过大或者热点,
  需要使用region sharding, 比如按照mid求余 % 和 / 10000商为key, 余数为offset
- List: 抽奖的奖池, 顶弹幕, 用于类似Stack Push/Pop操作
- Sortedset: 翻页, 排序, 有序集合, 杜绝range或者zrevrange返回的集合过大
- Hashs: 过小的时候会使用压缩列表, 过大的情况容易导致rehash内存浪费, 
  也杜绝返回hgetall, 对于小结构体, 建议直接使用memcache KV
- String: SET的EX/NX等KV拓展指令, SETNX可以用于分布式锁, SETEX聚合了SET + EXPIRE
- Sets: 类似Hashs, 无Value, 去重等
- 尽可能的PIPELINE指令, 但是避免集合过大
- 避免超大Value

### 分布式事务