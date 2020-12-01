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
- 野生goroutine如果panic 无法被recover， 需要构造一个 [func Go(x func())][4] 在其内部recover
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
- 利用已经封装好的标准库方法去消除代码中的err
- [将err封装进struct][3],利用struct的方法将if err隐藏起来，在方法内部判断struct.err
### Wrap errors
- 在之前的auth代码中错误会一直上抛，当程序顶端捕获到err的时候才会打印相关日志，但没有上下文信息。
- 在⬆️的基础上，我们改进上抛错误机制，引入`fmt.Errorf("xxxxx fail:%v",err)`,虽然有了具体错误信息，
但没有堆栈file:line信息，找起来麻烦，并且与sentinel err不兼容、等值判定失败。
- **you should only handle errors once.**
    - 如果发生错误，只允许一次日志。而不是每次err上抛都记录一次。
- 错误处理契约规定
    - 在出现err的情况下必须做处理、上抛，且不能对value做任何假设
    - 除非需要降级时候才可以不对err处理，并且将value设为default
- 日志记录只需要记录与错误相关且对调试有帮助的信息。
    - 错误必须被日志记录。
    - 应用程序处理错误，必须保证value的完整性。
    - 错误只被日志记录一次。
    - [github.com/pkg/errors][1]
- [pkg/errors][1]
    - [errors.Wrap][5] 附带错误的堆栈信息
    - errors.WithMessage 附带自定义信息
    - errors.Cause(err error) 返回最底层error
    - 通过 %+v 输出err的堆栈信息
    - 通过 errors.Wrap 不需要而外的log，并且可以携带上下文
- **如何使用 [pkg/errors][1]**
    - 在应用代码中通过 errors.New | errors.Errorf 返回错误
    - 如果是同程序内调用，直接`return err`,如果不透传而使用errors.Wrap,将会导致双倍的堆栈信息
    - 如果调用的是第三方库或者基础库，考虑使用errors.Wrap | errors.Wrapf 保存堆栈信息
    - 包内直接返回错误，而不是每个错误地方打日志
    - 在程序顶部或API入口处使用 %+v 打印并记录详细的堆栈信息
    - 通过使用 errors.Cause 获取 root error,再使用sentinel error判断
    - 作为基础库或第三方库，不应该返回errors.Wrap 而返回 root err, 只允许在业务代码中Wrap
    - 如果当前逻辑中不处理err，则可以将err wrap到上层，并携带足够的上下文，如 request params
    - 如果当前逻辑中已经处理err，则该err不再上抛 而是返回nil，如降级处理
## Go 1.13 errors
### Errors before Go 1.13
- 最简单的错误检查: err != nil
- 利用sentinel error: err == ErrNotFound
- 实现 error interface 的自定义 error struct 并利用断言: err.(*NotFoundError)
- 利用 fmt.Errorf 丢弃了除err的文本信息外的其它信息
### Errors in Go 1.13
- **Unwrap、Is、As**
- before 1.13: `fmt.Errorf("%v",err)` | 1.13: `fmt.Errorf("%w",err)` %w可以用于 Is As
- 可自定义 Is As 的错误方法
### errors vs github.com/pkg/errors
- errors 没有携带堆栈信息，pkg/errors 兼容 errors 的Is As
## Go 2 Error Inspection
- [https://go.googlesource.com/proposal/+/master/design/29934-error-values.md][2]

[1]: https://github.com/pkg/errors
[2]: https://go.googlesource.com/proposal/+/master/design/29934-error-values.md
[3]: https://github.com/XYZ0901/Go-000/blob/main/Week02/code/err_struct.go
[4]: https://github.com/XYZ0901/Go-000/blob/main/Week02/code/painc_goroutine.go
[5]: https://github.com/XYZ0901/Go-000/blob/main/Week02/code/err_wrap.go