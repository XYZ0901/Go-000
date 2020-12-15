# 学习笔记
## 作业
**在work目录下哦😉**
## 工程项目结构
- **在看本部分前先了解一下[Standard Go Project Layout][1]**
- 如果只是PoC或toy-product可以跳过该部分
- 当多人协作时则需要一个共用的目录结构,建议开发一个kit-tool,用于方便快速生成项目模板,统一目录布局.
- **/cmd**
    - 项目主干
    - 每个项目应该: /cmd/myapp/main.go 而不是: /cmd/myapp.go
    - 除非有必要,不添加额外代码.
- **/internal**
    - 私有程序、库代码, 只允许本项目引入和使用. 详情可以查看[Go1.4 release notes][2]
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
        - [上节课作业terryMao版][3]
- **Service Application Project - v1**
    - 某小破站的老项目布局 api,cmd,configs,internal 额外的还有 README.md,CHANGELOG,OWNERS
    - [internal][4]
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
- **Service Application Project - v2**
    - 某小破站的新式布局 api,cmd,configs,internal 额外的还有 README.md,CHANGELOG,OWNERS
    - internal
        - 为了避免同业务下有人跨目录引用内部的biz、data、service 等内部struct
        - biz
            - [业务逻辑层][6], 类似DDD的domain
            - 定义了业务逻辑实体, 业务实体应该在业务逻辑层, 定义了持久化接口
            - 在写业务逻辑的时候才知道数据应该如何被持久化, 持久化的interface应该被定义在业务逻辑层
        - data 
            - 类似DDD的repo, repo接口在这里定义, 使用依赖倒置原则
            - 业务数据访问层, 包括cache
            - 实现了biz定义的持久化接口逻辑
            - 事务暂时在这里实现
            - po(persistent Object) - 持久化对象, 与data层的数据结构一一对应
        - pkg - 实现业务逻辑的地方, 依赖倒置 service不依赖dao的具体struct 而是依赖dao的interface
        - service 
            - 实现了api定义的服务层, 类似DDD的application层
            - 实现dto -> do, 贫血模型
            - [IOC 控制反转、依赖注入][5] - 1、方便测试 2、单次初始化和复用
            - [https://github.com/google/wire][10]
            - 这里只应该有编排逻辑, 不应该有业务逻辑
    - 从根据功能组织到根据业务组织
    - LifeCycle
        - 手撸资源初始化与关闭-繁琐、易出错, 
        利用 [wire][10] 组织初始化代码, 非常方便快捷
## API设计
- **gRPC**
    - Week01
    - 方便、自动定义好接口协议
    - 三者合一 - pb、code、do
- **api管理(pb文件的存放)**
    - 由于每个项目的git不方便对外部暴露, 实现一个api的自动保存工程, 将所有repo中上传的pb转存到一个独立的api仓库, 
    不同部门的员工只能看到api仓库的pb文件, 不会看到别人的工程代码
    - 需要规范化检测, API Lint
    - 更方便code review
    - 权限管理, 目录中带有OWNERS
    - 可将api自动推送到各语言仓库, 员工可直接调用
- **api project layout**
    - 服务-应用-id-具体服务名称
    - options 自定义指令扩展
    - third_party 第三方pb调用
- **api Compatibility**
    - 向后兼容(非破坏性)的修改
        - 给API服务定义添加API接口
        - 给请求消息添加字段
        - 给响应消息添加字段
    - 向下不兼容(破坏性)的修改 - 如果需要重构 请更新大版本 如果v2
        - 删除或重命名服务、字段、方法或枚举值
        - 修改字段的类型
        - 修改现有请求的可见行为 - 业务逻辑有所变化
        - 给资源消息添加读/写字段
- **API Naming Conventions**
    - pb包名 - 应用id+版本号(`<package_name>.<version>`)
    - 生成的httpURL `/<package_name>.<version>.<service_name>/{method}`
- **API Primitive Fields**
    - 字段在 protobuf v2有required与optional分必选与可选, protobuf v3全部都是可选
    - Q: 如何在pb v3中区分可选与必选字段 - A: 将特定值重新定一个message结构体包含可分nil与传值
    - Protobuf作为强schema的描述文件, 也可以方便扩展
- **API Errors**
    - 一般http接口返回错误时均以200返回附与一个函数内错误, 这不可取、不便利、难以解析,
     建议直接返回http状态码
    - grpc也一样返回grpc的状态码或http2的状态码
    - 尽可能使用一组标准错误配合大量资源 - 例如sql找不到和redis找不到统一找不到(404)
    - 错误传播 
        - 不要透传, 用自己再生成的错误信息覆盖底层错误信息, 以便隐藏实现的详细信息和机密信息
    - 全局错误码 - 不建议使用
    - "小"错误 -> "大"错误
    - [app -> api/errors -> pb -> (1)enum & app -> krotos/errors -> errors.NotFound(1)][11]
    - kratos/errors -> status Error -> deep copy grpc -> client  
    - 不要把全局错误与grpc错误混合
- **API Design**
    - 单一字段更新接口过多 - 建议update接口只允许一个
    - 读写接口分离
    - [FieldMask(官方推荐)][12]

## 配置管理
- **环境变量(配置)**
    - Region、Zone、Cluster、Environment、Color、Discovery、AppID、Host等环境信息
- **静态配置**
    - 需要初始化的配置信息 http/grpc redis mysql等
    - 不鼓励on-the-fly 即不允许热更新 可以使用平滑更新
- **动态配置**
    - 尽量用基础类型配置可以考虑结合类似 [https://pkg.go.dev/expvar][13] 来使用
- **全局配置**
    - 统一依赖中间件
    - 配置中心全局模板
- **如何在初始化资源时传入配置信息**
    - ~~值传递~~
    - ~~全部选填传递~~
    - 指针传递 ✔️
    - 必选参数必填, 可选参数选填
    - [传入一系列函数指针改变参数数据][8]
    - [传入可返回前项配置的函数指针][9]
    - 混合模式 传入函数指针与传入参数混合写
- **Configuration**
    - 配置工具要素
        - 语义验证
        - 高亮
        - Lint
        - 格式化
    - 使用 YAML + Protobuf
    - config与配置文件解耦, 利用pb+wrapper+注解 可加验证规则 多语言统一配置文件
    - 避免复杂、多样配置、向简单化、向用户转变、必选与可选、
    配置的防御编程(是否合法)、权限与变更跟踪、配置版本与应用对齐、安全变更

## 包管理
- [Athens][14]
## 测试
- **测试金字塔**
    - 小型测试 - 单元测试 - kit库、中间件必须要
    - 中型测试 - 集成测试
    - 大型测试 - e2e测试 - 微服务建议只有这个
    - 混沌测试 - 中间件需要
- **单元测试基本要求**
    - 快速
    - 环境一致
    - 任意顺序
    - 并行
- **测试中初始化资源配置**
    - 利用sync.Once初始化资源配置,但不方便关闭
    - 改进 - 利用TestMain做初始化与关闭
    - 改进 - 利用subTest
- **TeseCase**
    - 表测试到子测试
    - 利用mock、fake
    - 利用依赖倒置
    - 利用docker compose构建
    - 利用yapi做接口测试
- **整体测试流程**
    1. 机遇git branch做feature开发
    2. 本地单元测试、自测
    3. 提交git merge
    4. ci自动验证单元测试
    5. 基于feature branch进行环境构建并功能测试
    6. 合并到master
    7. 再由CI跑一遍单元测试、用yapi做一次正常环境流程的集成测试
    8. 上线后进行回归测试
    
[1]: https://github.com/golang-standards/project-layout/blob/master/README_zh.md
[2]: https://golang.org/doc/go1.4#internalpackages
[3]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo1
[4]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo2
[5]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo3
[6]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo4
[7]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo5
[8]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo6
[9]: https://github.com/XYZ0901/Go-000/tree/main/Week04/demo7
[10]: https://github.com/google/wire
[11]: https://github.com/go-kratos/kratos/tree/v2/examples/kratos-demo
[12]: https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask
[13]: https://pkg.go.dev/expvar
[14]: https://github.com/gomods/athens