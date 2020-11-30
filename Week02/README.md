# 学习笔记
## Error vs Exception
### Exception In Other Language
- Error In C
    - 单返回值，通过指针入参获取数据，返回值int表示成功或失败
- Error In C++
    - 引入 exception 但无法知道被调用方抛出的是什么类型的异常
- Error In Java
    - 引入 checked exception 但不同的使用者会有不同处理方法，变得太司空见惯，严重程度只能人为区分，
    并且容易被使用者滥用，如经常 catch (e Exception) { // ignore }
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
- 只有真正意外、不可恢复的程序错误才会使用 panic , 如 索引越界、不可恢复的环境问题、栈溢出，才使用panic。除此之外都是error。
- go error 特点:
    - 简单
    - Plan for failure not success
    - 没有隐藏的控制流
    - 完全交给你来控制error
    - Error are values
## Error Type
### Sentinel Error
- sentinel error: 预定义错误,特定的不可能进行进一步处理的做法
- if err == ErrSomething { ... } 类似的sentinel error比如: io.EOF、syscall.ENOENT
- 最不灵活，必须利用==判断，无法提供上下文。只能利用error.Error()查看错误输出。
- 会变成API公共部分
    - 增加API表面积
    - 所有的接口都会被限制为只能返回该类型的错误，即使可以提供更具描述性的错误
- 在两个包中间产生依赖关系：无法二次修改现在包所返回的error，存在高耦合、无法重构
- **总结:尽可能避免sentinel errors**
### Error types
- Error type是实现了error接口的自定义类型，可以自定义需要的上下文及各种信息
- Error type是一个type 所以可以被**断言**用来获取更多的上下文信息
- VS Sentinel Error
    - Error type可以提供更多的上下文
    - 一样会public，与调用者产生强耦合，导致API变得脆弱。
    - 也需要尽量避免Error types
### Opaque errors (最标准、建议的方法)
- 只知道对或错 只能 err != nil
- 但无法携带上下文信息
- Assert errors for behaviour, not type
    - 通过定一个 interface ，然后暴露相关的Is方法去判断err，调用库内部断言。
- **具体选择哪种还是得需要看场景**
## Handing Error
### Indented flow is for errors
- err != nil 而不是 err == nil
### Eliminate error handling by eliminating errors
- 代码编写时可以直接返回err的 别用 err != nil
- 利用已经封装好的方法去消除代码中的err