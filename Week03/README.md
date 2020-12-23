### 学习笔记

1. Goroutine

    概念：
    1. Thread vs Process
    2. Concurrency vs Parallelism

    原则：
    1. Keep yourself busy or do the work yourself
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

    2. Never start a goroutine without knowing when it will stop
        a. 是否要通过 goroutine 后台执行应该由调用者决定；
        b. goroutine 的调用者负责处理整个生命周期, 确保退出(否则可能泄露)；

        ```
        func main() {
            go serveDebug()  // 不知道何时退出; 不推荐
            serveApp()
        }
        ```

        tips: log.Fatal 会导致程序直接退出，所以一般在init或者main中调用

    3. Leave concurrency to the caller


    4. Care About Incomplete Work (in goroutine)