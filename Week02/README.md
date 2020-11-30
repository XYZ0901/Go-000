# 学习笔记
## Error vs Exception
### Exception In Other Language
- Error In C
    - 单返回值，通过指针入参获取数据，返回值int表示成功或失败
- Error In C++
    - 引入 exception 但无法知道被调用方抛出的是什么类型的异常
- Error In Java
    - 
### Error In Go
- Go中的error只是一个普通的interface 包含一个 Error() string 方法
- 使用 errors.New() 创建一个error对象，返回的是 errorString 结构体的指针
- 利用指针 在基础库内部预定义了大量的 error, 用于返回及上层err与预定义的err做对比，
预防err文本内容一致但实际意义及环境不同的两个err对比成功。
- go支持多参数返回，一般最后一个参数是err，必须先判断err才使用value，除非你不关心value，即可忽略err.
- go的panic与别的语言的exception不一样，需谨慎或不使用，一般在api中第一个middleware就是recover panic.
- 野生goroutine如果panic 无法被recover， 需要构造一个 func Go(x func()) 在其内部recover
- 强依赖、配置文件: panic , 弱依赖: 不需要panic  
    - Q1: 案例: 如果数据库连不上但redis连得上,是否需要panic.
    - A1: 取决于业务，如果读多写少，可以先不panic，等待数据库重连。
    - Q2: 案例: 服务更新中导致gRPC初始化的client连不上
    - A2: 也是看业务，如果gRPC是Blocking(阻塞):等待重连、nonBlocking(非阻塞):立刻返回一个default、
    nonBlocking+timeout(非阻塞+超时/推荐):先尝试重连如果超时返回default
