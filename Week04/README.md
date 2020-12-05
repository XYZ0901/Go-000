# 学习笔记
## 作业
**在work目录下哦😉**
## 工程项目结构
- **在看本部分前先了解一下[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)**
- 如果只是PoC或toy-product可以跳过该部分
- 当多人协作时则需要一个共用的目录结构,建议开发一个kit-tool,用于方便快速生成项目模板,统一目录布局.
- **/cmd**
    - 项目主干
    - 每个项目应该: /cmd/myapp/main.go 而不是: /cmd/myapp.go
    - 除非有必要,不添加额外代码.
- **/internal**
    - 私有程序、库代码, 只允许本项目引入和使用. 详情可以查看[Go1.4 release notes](https://golang.org/doc/go1.4#internalpackages)
    - 针对每个项目都应该新建一个对应的目录, 而不是直接将.go文件放在本目录下.
    - 如果需要调用不暴露的公共函数, 可以在internal目录下添加pkg目录.
    - ~~如果只是单一项目, 可以考虑去掉项目目录, 直接将.go文件放本目录下.~~(不建议)
- **/pkg**
    - 可被外部程序调用的库代码, 会被其它项目引用, 所以放东西在里面时需要三思.
    - 该目录可参考go标准库的组织方式.(按功能划分目录)
    - /internal/pkg 用于本项目内调用, 不会被外部使用.
    - ~~相当于一个杂物间, 啥都往里放.~~
- **Kit Project Layout**
    - 每个公司都应该为不同的微服务建立一个统一的kit工具包项目(基础库/框架)和app项目.
    - 公司级建议只有一个, 如果有别的更好可以考虑合并, ~~或者通过行政手段干掉~~.
    - 不允许有vendor、不(少量)依赖第三方包. 让应用选择第三方包，而不是kit选择第三方.
    - 必要包需要依赖: grpc、proto
    - 可考虑封装插件或者fork代码的方式引入依赖
    - 特点: 统一、标准库方式布局(命名)、高度抽象、支持插件
- **/api**
    - API协议定义目录
    - 先把api文件安排进这个目录
- **/configs**
    - 配置文件(推荐yaml)
- **/test**
    - 较大的项目应该需要测试数据子目录, 如: /test/data 或 /test/testdata
    - 可以通过添加文件前缀`.`或`_`用于屏蔽go编译. 增加灵活性.
- ~~/src~~
    - ~~呵呵~~
- **Service Application Project**
    - 应用按照业务命名而不是部门命名.防止部门业务变更.
    - 多app方式: app目录内的每个微服务按照自己的全局唯一名称, 比如`account.service.vip`来建立目录,
    该名称还可以用于做服务发现.
    - app服务类型分类
        - **interface** 对外的BFF服务
        - **service** 提供对内的服务
        - **admin** 提供运营侧的微服务,允许更高权限,提供代码安全隔离. 这里与service共享数据, share pattern.
        - **job** 流式任务处理: 处理kafka、rabbitmq等消息队列的任务
        - **task** 定时任务, 类似cronjob, 部署到task托管平台中
        - cmd的本质:
            - 资源初始化、注销、监听、关闭
            - 初始化redis mysql dao service log 监听信号
        - [上节课作业terryMao版](https://github.com/XYZ0901/Go-000/tree/main/Week04/demo1)
- Service Application Project - v1
    - 某小破站的老项目布局 api,cmd,configs,internal 额外的还有 README.md,CHANGELOG,OWNERS
    - [internal](https://github.com/XYZ0901/Go-000/tree/main/Week04/demo2)
        - model - 各种结构体struct
        - dao - 访问mysql,redis等数据库方法 关联调用model, 面向的是一张表
        - service - 实现业务逻辑的地方, 依赖倒置 service不依赖dao的具体struct 而是依赖dao的interface
        - server - 依赖service, 放置grpc、http的起停、路由信息
    - 缺陷
        - 结构体会从model层层传到server层最后再经过json序列化.
        - model与表绑定 但有些字段需要屏蔽不从接口返回, 有些字段无法被json化, 需要转换.
        - 无法确定处理返回数据放在哪个位置
        - 补救措施
            - 引入DTO对象, 在有需求时对数据进行转换
    - 项目依赖路径: model -> dao -> service -> api(具有DTO转换)
        - 将cache数据方法从service全沉入dao层, 使得service层更专注业务, 从而cache miss放入dao
        - server层可被取消或替换掉
        - 不允许dto对象被dao引用
    - 整体按功能划分
        - 失血模型到贫血模型的转换
        - 失血模型: model层只存放数据结构,不实现任何逻辑
        - 贫血模型: model层为数据结构添加判断逻辑方法
        - ~~充血模型: 在贫血模型的基础上加入数据持久化的逻辑(不推荐)~~
- Service Application Project - v1
    - 某小破站的新式布局 api,cmd,configs,internal 额外的还有 README.md,CHANGELOG,OWNERS
    - internal
        - 为了避免同业务下有人跨目录引用内部的biz、data、service 等内部struct
        - biz
            - 业务逻辑层, 类似DDD的domain
            - 定义了业务逻辑实体, 业务实体应该在业务逻辑层, 定义了持久化接口
            - 在写业务逻辑的时候才知道数据应该如何被持久化, 持久化的interface应该被定义在业务逻辑层
        - data 
            - 类似DDD的repo, repo接口在这里定义, 使用依赖倒置原则
            - 业务数据访问层, 包括cache
            - 实现了biz定义的持久化接口逻辑
            - po(persistent Object) - 持久化对象, 与data层的数据结构一一对应
        - pkg - 实现业务逻辑的地方, 依赖倒置 service不依赖dao的具体struct 而是依赖dao的interface
        - service 
            - 实现了api定义的服务层, 类似DDD的application层
            - 实现dto -> do, 贫血模型
            - IOC 控制反转、依赖注入 - 1、方便测试 2、单次初始化和复用
            - [https://github.com/google/wire](https://github.com/google/wire)
            - 这里只应该有编排逻辑, 不应该有业务逻辑
    - 从根据功能组织到根据业务组织
    - LifeCycle
        - 手撸资源初始化与关闭-繁琐、易出错, 
        利用 [wire](https://blog.golang.org/wire) 组织初始化代码, 非常方便快捷
    

## API设计
## 配置管理
## 包管理
## 测试
## Reference