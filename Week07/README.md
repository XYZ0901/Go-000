# 学习笔记
## 作业
**在work目录下哦😉**
## 历史记录
### 功能模块
- 播发历史查看、播放进度同步, 离线型用户: app本地保留历史记录数据
- 可以考虑平台化, 在视频、文章、漫画等业务拓展接入
  - 变更功能: 添加记录、删除记录、清空历史
  - 读取功能: 按照timeline返回topN, 点查获取进度信息
  - 其它功能: 暂停/恢复记录, 首次观察增加经验等
- 具有极高tps写入, 高qps读取的业务服务, 分析清楚系统的hot path再投入优化 (可以考虑先本地再存储)
### 架构设计
- BFF(api-interface/history)
  - 数据组装
  - BFF的划分
    - 重要性
    - 垂直分层
    - 流量大小
    - 业务划分
  - 
- Service: history-service
  - 专注历史数据的持久化上
  - 负责数据的读、写、删、清理
  - 使用write-back的思路, 先写缓存再回写数据库
  - (在客户端上没有任何安全可言)
  - 写的核心逻辑
    - last-write win, 高频的用户端同步逻辑, 只需要最后一次数据持久化即可
    - 在in-process内存中, 定时定量聚合不同用户的同一个对象的最后一次进度
    - 使用kafka消息队列进行写入消峰
    - 为了同时保证用户数据可以被实时观察, 不能在上报进度后需要一阵子才能体现进度变化, 
      所以我们在内存中打包数据同时将数据实时写入到redis中, 这样保证了实时, 
      又避免了海量写入冲击存储
    - cache-aside读
    - kafka是为了高吞吐设计, 超高频的写入并不是最优, 所以内存聚合和分片算法比较重要, 
      按照uid来sharding数据, 使用region sharding, 打包一组数据当作一个kafka message
  - 写逻辑的数据流向: 实时写redis -> 内存维护用户数据 -> 定时/定量写入kafka
  - 读的核心逻辑: 历史数据, 实时写入redis后, 不会无限存储, 按流量截断, 分布式缓存中的数据不是完整数据,
    历史数据从redis sortedset中读取后, 如果发现尾部数据不足, 触发cache-aside 从hbase捞数据
- Job: history-job
  - job消费上游kafka的数据, 利用消息队列的堆积能力, 
    对于存储层的差速(消费能力跟不上生产速度), 可以进行一定的数据反压, 
    配合上游的service批量打包过来的数据持久化
  - 需要先读redis获取完整数据, 再batch write进hbase
  - **批量打包(pipeline)聚合数据, 将高频、密集的写请求write-back, 批量消费减少对存储的直接压力**
  - 上游 history-job 按照 uid region sharding聚合好的数据, 在job中消费取出, 为了节约传输过程, 
    以及history-service的in-process cache的内存使用, 我们只维护了用户uid和视频id列表, 最小化存储
    和传输, 所以需要额外从redis再读, 再持久化
  - HBase非常适合高密度写入
- Upstream: some-app some-api
  - 整个历史服务还会被一些外部grpc服务依赖, 所以history还充当了内网的grpc Provider
    
### 存储设计
- 数据库设计
  - HBase
  - 数据写入: PUT mid, value
  - 数据读取: 列表获取为 GET mid
    - redis cache miss的时候不会去查HBase
- 缓存设计
  - 每次产生的历史数据立马更新redis
  - 使用sorted set基于时间排序的列表, member为义务id
  - 数据读取: 历史页面使用sorted set排序查找 mget批量获取history_content内容
  - 点查进度: 直接查找history_content进行点查, 不再回源HBase
  - bitmap, bloom filter
  - 每次触发都有前置判定, 是否有更好的优化方案

### 可用性设计
- Write-Back
  - 在history-service中实时写入redis数据, 因此只需要重点优化缓存架构中, 扛住峰值的流量写入
  - 在内存中使用```map[int]map[int]struct{}```聚合数据, 利用chan在内部发送小消息, 
    在sendproc使用timer和定量判定logic发送到下游的kafka中
  - 在history-job中 获取kafka的数据后重新去redis回捞数据: history-content, 
    重构完整数据后存入HBase
    - 风险1: history-service重启过程中, 预聚合的消息丢失
    - 风险2: history-job读取redis构建数据, 但redis丢失 - redis需要有副本
  - 进行了trade-off 高收敛比的设计 将大流量数据聚合为小流量
- 聚合
  - 不把逻辑上移而只是数据上移, 在BFF层可以考虑实现Batch接口
  - 经过API Gateway的流量会高频触发per-rpc auth 给内网的identify-service带来了不少压力
    - 使用长链接, 一次验证不断使用
- 广播
  - 用户首次触发的行为, 需要发送消息给下游系统进行触发其它奖励, 如何减少这类一天只用一次的标记位缓存请求?
  - 使用in-process localcache ❌
  - 在写操作使用flag回传 ✅