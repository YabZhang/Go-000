### 学习笔记

1. Goroutine

    概念：
    1.1 Thread vs Process
    1.2 Concurrency vs Parallelism

    原则：
    1.1 Keep yourself busy or do the work yourself
        ```
            // 不鼓励的写法
            func main() {
                ...
                // goroutine 中执行业务逻辑
                go func() {
                    if err != http.ListenAndServe(":8080", nil); err != nil {
                        log.Fatal(err)
                    }
                }()
                select {} // 空 select 永远阻塞
            }

            // 比起空闲的等待，通常可以自己去执行任务
            func main() {
                http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                    fmt.Fprintln(w, "Hello, GopherCon SG")
                })

                if err := http.ListenAndServe(":8080", nil); err != nil {
                    log.Fatal(err) // os.Exit 直接退出，中断 defer 执行；不推荐
                }
            }
        ```

    1.2 Never start a goroutine without knowing when it will stop
        a. 是否要通过 goroutine 后台执行应该由调用者决定；
        b. goroutine 的调用者负责处理整个生命周期, 确保退出(否则可能泄露)；

        ```
        func main() {
            go serveDebug()  // 不知道何时退出; 不推荐
            serveApp()
        }
        ```

        tips: log.Fatal 会导致程序直接退出，所以一般在init或者main中调用

    1.3 Leave concurrency to the caller

    1.4 Care About Incomplete Work (in goroutine)

2. Memory model

    内存模型: https://golang.org/ref/mem，内存模型指定了跨 goroutine 的对公共变量并发读写操作必须满足的条件。
    同步语义: https://www.jianshu.com/p/5e44168f47a3

    2.1 Happens Before（先行发生）

        单个 goroutine 里的读写操作一定是按照程序指定的逻辑顺序来执行的。
        但是由于编译器和处理器的指令重排作用，不同 goroutine 观察到的执行顺序可能是不同的。

        在单独的goroutine中先行发生的顺序即是程序中表达的顺序。
        当下面条件满足时，对变量v的读操作r是 **被允许** 看到对v的写操作w的：

        1 r不先行发生于w
        2 在w后r前没有对v的其他写操作

        为了保证对变量v的读操作r看到对v的写操作w,要 **确保** w是r允许看到的唯一写操作。即当下面条件满足时，r 被保证看到w：

        1 w先行发生于r
        2 其他对共享变量v的写操作要么在w前，要么在r后。

    2.2 Synchronization

3. Package Sync

    Don't communicate by sharing memory, instead, share memory by communicating.

    Detecting Race Conditions With Go
        go build -race
        go test -race

        tips: 输出汇编 go tool compile -S {filename}; a++ 非原子

    Atomic CAS无锁编程
    Mutex 实现：
       * Barging: 释放者会去唤醒第一个等待的人(可能是自己)，提高吞吐量;(最高吞吐)
       * Handoff: 当锁释放时，会一直持有到第一个等待着准备好获取锁;(Park等待)
       * Spinning: 自旋在等待队列为空或者应用重度用锁时效果不错;(高频等待，容易导致饥饿)
       Go 1.8: Bargin + Spinning
       Go 1.9: 增加饥饿模式下停用自旋

    sync.Pool: 保存和复用临时对象，减少内存分配和GC压力；
        RingBuffer(定长FIFO) + 双向链表
        有可能被GC掉，所以不适合连接型对象管理


    errgroup
        WaitGroup + goroutine；并行工作流，优雅错误处理和终止，context传播，局部变量+闭包传递结果
        https://pkg.go.dev/golang.org/x/sync/errgroup


        https://github.com/go-kratos/kratos/tree/master/pkg/sync/errgroup
        * 野生 goroutine panic 捕获
        * 限定 errgroup 最大执行 goroutine 数，增加排队机制
        * context cancel 处理

4. channel

    Unbuffered channel（保证同步; Recive先于Send） vs buffered channel
    Latencies due to under-sized buffer

    Go Concurrency Patterns: 参考官网文章

5. package Context

    Go1.7 开始引入context 包，它使得跨API边界的请求范围元数据、取消信号和截止信号很容易传递给处理请求所涉及到其他 goroutine (显示传递).
    其他语言: Thread Local Storage, XXXContext

    两种方式集成Context:
    1. func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error)
    2. func (r *Request) WithContext(ctx context.Context) *Request

    Do not store Contexts inside a struct type! (执行态的上下文环境，尽量只传递如函数签名)

    Debugging or tracing data is safe to pass in a Context (为了防止数据竞争，Context中的数据应该只读)
    Context.Value should inform, not control;

    When a Context is canceled, all Contexts derived from that are also canceled. => close(channel)
    All blocking/long operation should be cancelable.

    https://talk.golang.org/2014/gotham-context.slide#1


### 作业题目

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。