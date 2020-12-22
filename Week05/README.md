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