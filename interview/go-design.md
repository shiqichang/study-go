## Go 语言设计与实现

### 一、编译原理

抽象语法书 (AST Abstract Syntax Tree)

静态单赋值 (SSA static single assignment)

指令集：

- 复杂指令集计算机 CISC
- 精简指令集计算机 RISC

编译器的4个阶段：

1. 词法与语法分析：
   1. 词法分析：解析源代码文件，将字符串序列转换成 Token 序列；
   2. 词法解析器 (lexer)
   3. 语法分析：将词法分析生成的 Token 按编程语言定义好的文法 (Grammar) 自下而上或自上而下地规约，每个 Go 源代码文件最终会被归纳成一个 SourceFile 结构；
2. 类型检查和 AST 转换：
   1. 类型检查：
      1. 常量、类型和函数名及类型；
      2. 变量的赋值和初始化；
      3. 函数和闭包的主体；
      4. 哈希键值对的类型；
      5. 导入函数体；
      6. 外部的声明；
3. 通用 SSA 生成
4. 最后的机器码生成

抽象语法树会经历类型检查、SSA 中间代码生成、机器码生成三个阶段：

1. 检查常量、类型和函数的类型；
2. 处理变量的赋值；
3. 对函数的主体进行类型检查；
4. 决定如何捕获变量；
5. 检查内联函数的类型；
6. 进行逃逸分析；
7. 将闭包的主体转换成引用的捕获变量；
8. 编译顶层函数；
9. 检查外部依赖的声明；

顶层声明 5 大类型：常量、类型、变量、函数、方法；

### 二、数据结构

#### 1. 数组

对**数组**的访问和赋值需要同时依赖编译器和运行时，它的大多数操作在编译期间都会转换成直接读写内存，在中间代码生成期间，编译器还会插入运行时方法 `runtime.panicIndex` 调用防止越界错误。

```go
package main

var arr1 = [3]int{1, 2, 3}
var arr2 = [...]int{1, 2, 3}
```

访问数组的索引不能是非整数、负数。

#### 2. 切片

##### 数据结构

```go
package main

type SliceHeader struct {
    Data uintptr // 指向数组的指针
    Len int      // 当前切片的长度
    Cap int      // 当前切片的容量，即 Data 数组的大小
}
```

##### 初始化

```go
package main

// arr[0:3] or slice[0:3]   // 通过下标的方式获取数组或切片的一部分
var slice = []int{1, 2, 3}  // 使用字面量初始化新的切片
var slice = make([]int, 10) // 使用关键字 make 创建切片
```

#### 3. 哈希表

##### 解决冲突

```go
package main

var index = hash("key") % array.len
```

1. 开放寻址法：依次探测和比较数组中的元素以判断目标键值对是否存在于哈希表中，实现哈希表底层的数据结构是数组；
2. 拉链法：使用数组加链表，底层数据结构是链表数组
   1. 找到键相同的键值对 - 更新键对应的值；
   2. 没有找到键相同的键值对 - 在链表的末尾追加新的键值对；
   3. 计算哈希、定位桶和遍历链表是哈希表读写操作的主要开销；
   4. 装载因子 := 元素数量 / 桶数量

哈希在每个桶中存储键对应哈希的前 8 位，每个桶中都只能存储 8 个键值对。

#### 4. 字符串

```go
package main

type StringHeader struct {
    Data uintptr // 指向字节数组的指针
    Len  int     // 数组的大小
}
```

```go
package main

var str1 = "this is a string"
var str2 = `this is another
string`
```

### 三、语言基础

#### 1. 函数调用

在 x86_64 机器上使用 C 语言调用函数时，参数都是通过寄存器和栈传递的，其中：

1. 六个及六个以下的参数会按照顺序分别使用 edi、esi、edx、ecx、r8d 和 r9d 六个寄存器传递；
2. 六个以上的参数会使用栈传递，函数的参数会以从右到左的顺序依次存入栈中；

Go 语言使用栈传递参数和返回值。

##### 参数传递

Go 语言选择了**传值**的方式。无论是传递基本类型、结构体还是指针，都会对传递的参数进行拷贝。

- 通过堆栈传递参数，入栈的顺序是从右到左，参数的计算是从左到右；
- 函数返回值通过堆栈传递，并由调用者预先分配内存空间；
- 调用函数时都是传值，接收方会对入参进行复制并计算。

#### 2. 接口

Go 语言中接口的实现都是隐式的。

Go 语言的接口类型不是任意类型。

##### 指针和接口

|                      | 结构体实现接口 | 结构体指针实现接口 |
| -------------------- | -------------- | ------------------ |
| 结构体初始化变量     | 通过           | 不通过             |
| 结构体指针初始化变量 | 通过           | 通过               |

##### 数据结构

根据接口类型是否包含一组方法，将接口类型分为：

- 使用 `runtime.iface` 表示包含方法的接口；
- 使用 `runtime.eface` 表示不包含任何方法的 `interface{}` 类型；

```go
package main

import "unsafe"

type eface struct { // 16 字节
    _type *_type         // 类型结构体
    data  unsafe.Pointer
}
```

`interface{}` 类型不包含任何方法，只包含指向底层数据和类型的两个指针。Go 语言的任意类型都可以转换成 `interface{}`.

```go
package main

import "unsafe"

type iface struct { // 16 字节
    tab  *itab          // itab 结构体
    data unsafe.Pointer
}
```

类型转换、类型断言、动态派发（Dynamic dispatch）

#### 3. 反射

两个函数：

- `reflect.TypeOf` 能获取类型信息；
- `reflect.ValueOf` 能获取数据的运行时表示；

两个类型：`reflect.Type` 和 `reflect.Value`.

```go
package main

type Type interface {
    Align() int
    FieldAlign() int
    Method(int) Method
    MethodByName(string) (Method, bool) // 获取当前类型对应方法的引用
    NumMethod() int
    Implements(u Type) bool // 判断当前类型是否包含某个接口
}

type Value struct { 
    // 包含过滤的或者未导出的字段
}

func (v Value) Addr() Value
func (v Value) Bool() bool
func (v Value) Bytes() []byte
```

运行时反射是程序在运行期间检查自身结构的一种方式。

##### 三大法则

- 从 `interface{}` 变量可以反射出反射对象；
- 从反射对象可以获取 `interface{}` 变量；
- 要修改反射对象，其值必须可设置。

`Interface Value` -- `TypeOf, ValueOf` -> `Reflection Object`

`Reflection Object` -- `Interface` -> `Interface Value`

### 四、常用关键字

#### 1. for 和 range

Go 遍历数组和切片时，会复用变量；

哈希表的随机遍历原理。

#### 2. select

I/O 多路复用：`select`, `poll`, `epoll`

- `select` 能在 `Channel` 上进行非阻塞的收发操作；
- `select` 在遇到多个 `Channel` 同时响应时，会随机执行一种情况；

`Type must be a pointer, channel, func, interface, map, or slice type.`

1. 空的 `select` 语句会被转换成调用 `runtime.block` 直接挂起当前 `goroutine`;
2. 若 `select` 中只有一个 `case`, 编译器会将其转换成 `if ch == nil { block }; n;` 表达式：首先判断操作的 `channel` 是不是空的，然后执行 `case` 与剧中的内容；
3. 若 `select` 只包含两个 `case` 且其中一个是 `default`, 那么会使用 `runtime.selectnbrecv` 和 `runtime.selectnbsend` 非阻塞地执行收发操作；
4. 在默认情况下会通过 `runtime.selectgo` 获取执行 `case` 的索引，并通过多个 `if` 语句执行对应 `case` 中的代码。

#### 3. defer

`Defer` 会在当前函数返回前执行传入的函数。经常被用于关闭文件描述符、关闭数据连接以及解锁资源。

调用 `defer` 关键字会使用值传递。向 `defer` 关键字传入匿名函数时，拷贝的是函数指针。

处理 `defer` 关键字的三种方法：堆分配、栈分配和开放编码。

`defer` 关键字的插入顺序是从后向前，执行顺序是从前向后。

- 后调用的 `defer` 函数会先执行：
  - 后调用的 `defer` 函数会被追加到 Goroutine `_defer` 链表的最前面
  - 运行 `runtime._defer` 是从前到后执行的；
- 函数的参数会被预先计算：
  - 调用 `runtime.deferproc` 函数创建新的延迟调用时会立刻拷贝函数的参数，函数的参数不会等到真正执行时计算。

#### 4. panic 和 recover

- `panic` 能改变程序的控制流，调用 `panic` 会立刻停止执行当前函数的剩余代码，并在当前 Goroutine 中递归执行调用方的 `defer`;
- `recover` 可以终止 `panic` 造成的程序崩溃。它是一个只能在 `defer` 中发挥作用的函数，在其他作用域中调用不会发挥作用；

现象：

- `panic` 只会触发当前 Goroutine 的 `defer`;
- `recover` 只有在 `defer` 中调用才会生效；
- `panic` 允许在 `defer` 中嵌套多次调用。

#### 5. make 和 new

- `make` 的作用是初始化内置的数据结构，如切片、哈希和 `Channel`;
- `new` 的作用是根据传入的类型分配一片内存空间，并返回指向这片内存空间的指针。

### 五、并发编程

#### 1. 上下文 Context

上下文 `context.Context` 用来设置截止日期、同步信号，传递请求相关值的结构体。

在 Goroutine 构成的树形结构中，对信号进行同步以减少计算资源的浪费是 `context.Context` 最大的作用。Go 服务的每一个请求都是通过单独的 Goroutine 处理的，HTTP/RPC 请求的处理器会启动新的 Goroutine 访问数据库和其他服务。

```go
package main

import "time"

type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
    return
}

func (*emptyCtx) Done() <-chan struct{} {
    return nil
}

func (*emptyCtx) Err() error {
    return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
    return nil
}
```

#### 2. 同步原语与锁

Go 是一个原生支持用户态进程 (Goroutine) 的语言。锁是一种并发编程中的同步原语 (Synchronization Primitives)，它能保证多个 Goroutine 在访问同一片内存时不会出现竞争条件 (Race condition) 等问题。

##### 基本原语

`sync` 包中的一些基本原语：`sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, `sync.Cond`.

###### Mutex

```go
package main

type Mutex struct {
    state int32  // 当前互斥锁的状态
    sema  uint32 // 用于控制锁状态的信号量
}
```

互斥锁：只占8个字节的结构体

- `mutexLocked`: 互斥锁的锁定状态；
- `mutexWoken`: 从正常模式被唤醒；
- `mutexStarving`: 进入饥饿状态；
- `waitersCount`: 等待的 Goroutine 个数。

`sync.Mutex` 有两种模式：正常模式和饥饿模式。饥饿模式是为了保证公平性。

- 加锁：`sync.Mutex.Lock`
- 解锁：`sync.Mutex.Unlock`

`sync.Mutex.lockSlow` 尝试通过自旋 (Spinning) 等方式等待锁的释放。

自旋是一种多线程同步机制。当前的进程进入自旋的过程中会一直保持 CPU 的占用，持续检查某个条件是否为真。在多核 CPU 上，自旋可以避免 `Goroutine` 的切换。

`Goroutine` 进入自旋的条件：

- 互斥锁只有在普通模式下才能进入自旋；
- `runtime.sync_runtime_canSpin` 需要返回 `true`:
  - 运行在多 CPU 的机器上；
  - 当前 `Goroutine` 为了获取该锁进入自旋的次数小于4；
  - 当前机器上至少存在一个正在运行的处理器 P 且处理的运行队列为空。

###### RWMutex

读写互斥锁 `sync.RWMutex` 是细粒度的互斥锁。它不限制资源的并发读，但读写、写写操作无法并行执行。

- 写操作使用 `sync.RWMutex.Lock` 和 `sync.RWMutex.Unlock` 方法；
- 读操作使用 `sync.RWMutex.RLock` 和 `sync.RWMutex.RUnlock` 方法。

###### WaitGroup

它可以等待一组 Goroutine 的返回。使用场景：批量发出 RPC 或 HTTP 请求。

- `sync.WaitGroup` 必须在 `sync.WaitGroup.Wait` 方法返回之后才能被重新使用；
- `sync.WaitGroup.Done` 方法只是对 `sync.WaitGroup.Add` 方法的简单封装，向 `sync.WaitGroup.Add` 方法传入任意负数（需要保证计数器非负），快速将计数器归零，以唤醒等待的 Goroutine;
- 可以同时有多个 Goroutine 等待当前 `sync.WaitGroup` 计数器归零，这些 Goroutine 会被同时唤醒。

###### Once

保证在 Go 程序运行期间的某段代码只会执行一次。

`sync.Once.Do` 方法会接收一个入参为空的函数：

- 若传入的函数已经执行过，会直接返回；
- 若传入的函数没有执行过，会调用 `sync.Once.doSlow` 执行传入的函数。

注意：

- `sync.Once.Do` 方法中传入的函数只会被执行一次，哪怕函数中发生了 `panic`;
- 两次调用 `sync.Once.Do` 方法传入不同的函数，只会执行第一次调用传入的函数。

###### Cond

它可以让一组 Goroutine 都在满足特定条件时被唤醒。

##### 扩展原语

###### ErrGroup

在一组 `Goroutine` 中，提供了同步、错误传播、上下文取消的功能。

```go
package main

import "sync"

type Group struct {
    cancel  func()
    wg      sync.WaitGroup
    errOnce sync.Once
    err     error
}
```

- `cancel`: 创建 `context.Context` 时返回的取消函数，用于在多个 `Goroutine` 中同步取消信号；
- `wg`: 用于等待一组 `Goroutine` 完成子任务的同步原语；
- `errOnce`: 用于保证只接收一个子任务返回的错误。

###### Semaphore

信号量是在并发编程中常见的一种同步机制。控制访问资源的进程数量时会用到信号量，它会保证持有的计数器在 0 和初始化的权重之间波动。

###### SingleFlight

它能够在一个服务中抑制对下游的多次请求。

- 缓存穿透
  - 缓存和数据库中都没有的数据，而用户不断发起请求。
  - 解决方案：
    - 接口层增加校验，如用户权限校验、id做基础校验，id<=0的直接拦截；
    - 将 key-value 对写成 **key-null**，缓存有效时间可以设置短点。
- 缓存击穿
  - 缓存中没有但数据库中有的数据，一般是缓存时间到期，而并发用户特别多。
  - 解决方案：
    - 设置热点数据永远不过期；
    - 接口限流与熔断，降级；
    - 布隆过滤器 bloomfilter, 类似于一个 hashset, 用于快速判断某个元素是否存在于集合中；
    - 加互斥锁。
- 缓存雪崩
  - 缓存中数据大批量到过期时间，而查询数据量巨大，引起数据库压力过大甚至down机。
  - 解决方案：
    - 缓存数据的过期时间设置随机；
    - 若缓存数据库是分布式部署，将热点数据均匀分布在不同的缓存数据库中；
    - 设置热点数据永远不过期。

#### 3. 计时器

##### 设计原理

###### 全局四叉堆

Go 1.10 之前的计时器都使用最小四叉堆实现。

###### 分片四叉堆

Go 1.10 将全局的四叉堆分割成了 64 个更小的四叉堆。

###### 网络轮询器

所有的计时器都以最小四叉堆的形式存储在处理器中。

##### 数据结构

```go
package main

type timer struct {
    pp puintptr
  
    when     int64 // 当前计时器被唤醒的时间
    period   int64 // 两次被唤醒的间隔
    f        func(interface{}, uintptr) // 每次被唤醒时会调用的函数
    arg      interface{} // 调用 f 传入的参数
    seq      uintptr
    nextwhen int64 // 计时器处于 timerModifiedXX 状态时，用于设置 when 字段
    status   uint32 // 计数器的状态
}
```

`time.Timer` 计时器必须通过 `time.NewTimer`、`time.AfterFunc` 或 `time.After` 创建。计时器失效时，订阅计时器 `Channel` 的 `Goroutine` 会收到计时器失效的时间。

##### 状态机

包含以下 10 种：

- timerNoStatus: 还没有设置状态，计时器不在堆上
- timerWaiting: 等待触发，计时器在处理器的堆上
- timerRunning: 运行计时器函数，停留时间较短，计时器在处理器的堆上
- timerDeleted: 被删除，计时器在处理器的堆上
- timerRemoving: 正在被删除，停留时间较短，计时器在处理器的堆上
- timerRemoved: 已经被停止并从堆中删除，计时器不在堆上
- timerModifying: 正在被修改，停留时间较短，计时器在处理器的堆上
- timerModifiedEarlier: 被修改到了更早的时间，计时器在处理器的堆上，可能位于错误的位置上，需要重新排序
- timerModifiedLater: 被修改到了更晚的时间，计时器在处理器的堆上，可能位于错误的位置上，需要重新排序
- timerMoving; 已经被修改正在被移动，停留时间较短，计时器在处理器的堆上

##### 触发计时器

###### 调度器

用来运行处理器中计时器的函数。检查器中的计时器是否准备就绪。

###### 系统监控

检查是否有未执行的到期计时器。

#### 4. Channel

Go 语言中最常见的设计模式：不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存。

##### 设计原理

###### 先进先出

- 先从 Channel 读取数据的 Goroutine 会先接收到数据；
- 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利。

###### 无锁管道

无锁 lock-free

- 悲观并发控制 Pessimistic
- 乐观并发控制 Optimistic

##### 创建管道

`make(chan int, 10)`

##### 发送数据

`ch <- i`

- 直接发送
- 缓冲区
- 阻塞发送

##### 接收数据

`i <- ch`, `i, ok <- ch`

- 直接接收
- 缓冲区
- 阻塞接收

##### 关闭管道

#### 5. 调度器

##### 设计原理

- 单线程调度器
- 多线程调度器
- 任务窃取调度器
- 抢占式调度器
  - 基于协作的抢占式调度器
  - 基于信号的抢占式调度器
- 非均匀内存访问调度器

##### 数据结构

- G - Goroutine, 一个待执行的任务
- M - 操作系统的线程，由操作系统的调度器调度和管理
- P - 处理器，运行在线程上的本地调度器

#### 6. 网络轮询器

网络轮询器是对 I/O 多路复用的封装。

- 初始化
- 轮询事件
- 事件循环
- 截止日期

#### 7. 系统监控

监控循环

- 检查死锁
- 运行计时器
- 轮询网络
- 抢占处理器
- 垃圾回收

### 六、内存管理

#### 1. 内存分配器

内存空间两个重要区域：栈区(Stack) 和 堆区(Heap).

函数调用的参数、返回值及局部变量大都被分配在栈上，这部分内存由编译器进行管理。

堆中的对象由内存分配器分配并由垃圾收集器回收。

内存管理包含三个组件：用户程序 (Mutator)、分配器(Allocator)、收集器 (Collector).

分配方法：

- 线性分配器 (Bump Allocator)
- 空闲链表分配器 (Free-List Allocator)

#### 2. 垃圾收集器

- 标记清除
- 三色抽象
- 屏障技术
- 增量和并发

#### 3. 栈内存管理

栈操作

- 栈初始化
- 栈分配
- 栈扩容
- 栈缩容

### 七、元编程

#### 1. 插件系统

#### 2. 代码生成

### 八、标准库

#### 1. JSON

#### 2. HTTP

#### 3. 数据库

