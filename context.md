接口 `context.Context` 的主要作用：在 Go 进程中传递信号。

```go
package main

import "time"

type Context interface {
    Dealine() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

- `Deadline`: 返回 `Context` 被取消的时间，即生命周期结束的时间；
- `Done`: 返回一个 `channel`, 它会在 `Context` 生命周期结束之后被关闭，多次调用返回的是同一个 `channel`;
- `Err`: 返回 `Context` 结束的原因，只会在 `Done` 方法返回的 `channel` 关闭之后才返回非空值；
- `Value`: 从 `Context` 中获取对应 `Key` 的值。

通过 `Context` 来串联所有的 `goroutine`, 所有的 `goroutine` 根据 `Context` 的状态来决定是否要继续执行。

创建一个 Context

1. `context.emptyCtx`;
2. `context.Background()` 表示最顶层的 `Context`, 其他 `Context` 都应该由 `Background` 衍生而来；

2. `context.TODO()` 用于还不确定使用哪个 `Context`;

3. 如果没有通过参数接收到 `Context`, 就会使用 `Background` 作为初始的 `Context` 向后传递。

衍生 `Context` 

1. `context.cancelCtx`;
2. 通过 4 个方法来完成：
   1. `WithCancel`: 会创建一个 `cancelCtx` 
   2. `WithDeadline`: 会创建一个 `timerCtx`
   3. `WithTimeout`: 会直接调用 `WithDeadline`
   4. `WithValue`: 会创建一个 `valueCtx`

`Context` 是线程安全的，可以在多个 `goroutine` 之间使用。

