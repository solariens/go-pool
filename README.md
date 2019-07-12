# go-pool

>  一个 goroutine 协程池的实现 ，所有的操作都是预先分配好可以执行的goroutine,避免在连续的实际任务执行中启动新的协程所带来的损耗

>  a goroutine pool implements

>  go-pool是由三部分组成：job、dispatcher、worker
>  job为worker执行的任务，每一个任务都必须实现Do()、Name()和Error()三个方法，Do()方法是任务的具体执行逻辑，Name()方法返回此次任务的名称，Error()方法负责接收错误信息
>  dispatcher维护了一个全局队列，负责接收任务并存入队列中，并预先生成指定数量的worker，从队列中取出对应的任务并分发到空闲的worker中，如果没有空闲的worker则随机选择一个worker分发任务
>  worker维护任务列表，生成goroutine监听任务的到来，并执行任务的Do方法运行任务

    go test -bench .
    goos: darwin
    goarch: amd64
    pkg: github.com/solariens/go-pool
    Benchmark_work_10_10-4           1000000              1484 ns/op
    Benchmark_work_10_100-4          2000000               827 ns/op
    Benchmark_work_100_100-4         1000000              1155 ns/op
    Benchmark_work_100_1000-4        2000000               776 ns/op
    Benchmark_work_100_10000-4       2000000               721 ns/op
    Benchmark_work_1000_10000-4      2000000               798 ns/op
    PASS
    ok      github.com/solariens/go-pool  13.930s