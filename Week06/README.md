# 学习笔记

## 案例 - 评论系统架构设计
- Q: 何为平民架构
- A: 在基础设施~~一坨屎~~尚未搭建完整的情况下, 有一个特别良好的架构.

### 功能模块

### 架构设计
- API -> BFF
- BFF(Comment)
    - 作为服务编排需要依赖三方服务, account, filter ...
    - 提供可用性
- comment-service
    - 服务层, 去平台业务的逻辑, 专注于API的实现
    - 读
        - Cache-Aside模式
            - 先读缓存(redis), 再读存储(mysql)
            - cache rebuild全部数据,导致不合理难处理, 利用read ahead(预读)的思路处理服务
        - 回源
            - 将cache miss的信息通知到comment-job, 让comment-job去做rebuild cache, 
            从db获取缓存并更新到redis, 从而解决thundering herd现象
    - 写
        - 将写逻辑传入mq(kafka), 利用comment-job消费kafka写入db更新redis 从而消峰
        - 利用hash(comment_subject) % N(partitions)将数据分发到kafka的多个partition从而使得全局并行, 局部串行
- comment-job
    - 利用mq(kafka)做消峰处理
    - 先写db, 再写redis
- comment-admin
    - 运营与管理能力, 从业务中独立出来
    - 与service共享数据与存储
    - 利用Canal(中间件)订阅binlog的数据解析成es(ElasticSearch)的语句写入es, 可以添加joiner去合并别表
    - 千万不要把mysql作为一个分析性数据库使用
- ps.
    - 架构设计等同于数据设计, 梳理清楚数据的走向与逻辑
    - 避免环形依赖, 数据双向请求
    
### 存储设计
- 表设计具体看ppt和视频
- tips
    - 存储类型尽量小
    - 利用bits做多属性状态
    - 利用root, parent做层级
    - 因为有sharding 所以一般不需要join
    - 将一个大事务(三个表一起更新)拆为小事务(content表不用事务,另外两个用事务)
    - pc端使用预读, 事先缓存多页数据
    - 移动端使用瀑布流, 利用last ID请求后续数据而不是PageNum
   
### 缓存设计
- 三个缓存表 具体看ppt和视频
- tip
    - score排序从 comment_id -> 楼层号 占用更小内存
    - 不适合大查询
    - sorted set在追加数据前先续命再添加 防止数据丢失

### 可用性设计 - SingleFlight
- 利用SingleFlight和comment-service的定时local cache防止大量的读穿透与回源
- 不建议使用分布式锁
- [SingleFlight](https://pkg.go.dev/golang.org/x/sync/singleflight)

### 可用性设计 - 热点
- 提升吞吐
- 拆分主干
- 利用环形数组统计热点