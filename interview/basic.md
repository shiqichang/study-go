## 基础部分

### golang 中 make 和 new 的区别

- make 和 new 都是 golang 用来分配内存的内建函数，且在堆上分配内存。make 既分配内存，也初始化内存。new 只将内存清零，并没有初始化内存；
- make 返回的是引用类型本身，new 返回的是指向类型的指针；
- make 只能用来分配及初始化类型为 slice/map/channel 的数据；new 可以分配任意类型的数据。

### for range 时地址会发生改变吗

不会

### go defer, 多个 defer 的顺序，defer 在什么时机会修改返回值

defer 延迟函数，释放资源，收尾工作；如释放锁、关闭文件、关闭连接、捕获 panic
defer 函数紧跟在资源打开后面，否则 defer 可能得不到执行，导致内存泄露
多个 defer 调用顺序：FILO，即压入栈
defer、return、return value 执行顺序：return -> return value -> defer. defer 可以修改函数最终返回值
修改时机：有名返回值或者函数返回指针

```go
// 有名返回值
func b() (i int) {
	defer func() {
		i++
		fmt.Println("defer2:", i)
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	return
}

// 函数返回指针
func c() *int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2:", i)
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	return &i
}
```

### rune 类型

golang 中的字符串底层实现是通过 byte 数组。中文字符在 unicode 下占 2 个字节，在 utf-8 下占 3 个字节。golang 默认编码是 utf-8。
byte 等同于 int8：用来处理 ascii 字符；
rune 等同于 int32：用于来处理 unicode 或 utf-8 字符。

```go
var str = "hello 你好"
fmt.Println("len(str)", len(str)) // 12
fmt.Println("RuneCountInString:", utf8.RuneCountInString(str)) // 8
fmt.Println("rune:", len([]rune(str))) // 8
```

### golang 中解析 tag 是怎么实现的，反射原理是什么

gorm, json, yaml, gRPC, protobuf, gin.Bind() 都是通过反射实现的。

### golang 函数传入结构体时，传值还是传指针

golang 的函数参数传递都是值传递。

### golang 的 slice 底层数据结构和特性

底层数据结构：由一个 array 指针指向底层数组。

### golang 的 select 底层数据结构和特性

为 golang 提供多路 I/O 复用机制。
底层数据结构：select 语句和执行函数。

```go
select {
	case <- chan1:
		// chan1 成功读到数据
	case chan2 <- 1:
		// 成功向 chan2 写入数据
	default:
		// 以上均未成功
}
```

- select 操作至少要有一个 case 语句，出现读写 nil 的 channel 该分支会忽略，在 nil 的 channel 上操作会报错；
- select 仅支持管道，且是单协程操作；
- 每个 case 语句仅能处理一个管道，要么读要么写；
- 多个 case 的执行顺序是随机的；
- 存在 default 语句，select 将不会阻塞，但会影响性能。

### golang 的 defer 底层数据结构和特性

每个 defer 语句对应一个 _defer 实例，多个实例使用指针连接起来形成一个单链表，保存在 goroutine 数据结构中。

### 单引号、双引号、反引号的区别

- 单引号：表示 byte 类型或 rune 类型，对应 uint8 和 int32，默认是 rune 类型。byte 用来强调数据是 raw data，不是数字。rune 用来表示 Unicode 的 code point.
- 双引号：字符串，实际上是字符数组。
- 反引号：字符串字面量，不支持任何转义序列。

## map 相关

### 注意事项，是否并发安全

一定要初始化，否则 panic；
go 语言内建的 map 对象不是线程安全的，并行读写时运行时会检查，遇到并发问题会报错。

### map 中删除一个 key，内存会释放吗

在很大的 map 中，delete 操作没有真正释放内存，可能导致内存 OOM。一般做法是重建 map。go-zero 中内置了 safemap 的容器组件。
golang 释放 map 内存的方法：首先删除 map 中所有 key，map 占用内存仍处于【使用状态】；然后 map 置为 nil，map 占用内存处于【空闲状态】；最后处于空闲状态内存，一定时间内在下次申请可重复被使用，不必再向操作系统申请。

### 怎么处理对 map 的并发访问

方法一：使用内置 sync.Map
方法二：使用读写锁实现并发安全 map

### map 的数据结构是什么，怎么实现扩容

golang 中 map 是一个 kv 对集合。底层使用 hash table，用链表来解决冲突，出现冲突时，不是每一个 key 都申请一个结构通过链表串起来，而是以 bmap 为最小粒度挂载，一个 bmap 可以放 8 个 kv。在哈希函数的选择上，会在程序启动时，检测 cpu 是否支持 aes，如果支持，则使用 aes hash，否则使用 memhash。每个 map 的底层结构是 hmap，是有若干个结构为 bmap 的 bucket 组成的数组。每个 bucket 底层都采用链表结构。

### slices 能作为 map 类型的 key 吗

在 golang 规范中，可比较的类型都可以作为 map key，包括：
boolean 布尔值
numeric 数字	包括整型、浮点型，以及复数
string 字符串
pointer 指针	两个指针类型相等，表示两指针指向同一个变量或者同为nil
channel 通道	两个通道类型相等，表示两个通道是被相同的make调用创建的或者同为nil
interface 接口	两个接口类型相等，表示两个接口类型 的动态类型 和 动态值都相等 或者 两接口类型 同为 nil
structs、arrays	只包含以上类型元素

不能作为map key 的类型包括：
slices
maps
functions

## context 相关

### context 的结构，使用场景，用途

Go 的 Context 的数据结构包含 Deadline，Done，Err，Value，Deadline 方法返回一个 time.Time，表示当前 Context 应该结束的时间，ok 则表示有结束时间，Done 方法当 Context 被取消或者超时时候返回的一个 close 的 channel，告诉给 context 相关的函数要停止当前工作然后返回了，Err 表示 context 被取消的原因，Value 方法表示 context 实现共享数据存储的地方，是协程安全的。context 在业务中是经常被使用的。

context 的使用:
对于 goroutine，他们的创建和调用关系总是像层层调用进行的，就像一个树状结构，而更靠顶部的 context 应该有办法主动关闭下属的 goroutine 的执行。为了实现这种关系，context 也是一个树状结构，叶子节点总是由根节点衍生出来的。

- 要创建 context 树，第一步应该得到根节点，context.Background 函数的返回值就是根节点。该 context 一般由接收请求的第一个 goroutine 创建。它不能被取消，也没有值，也没有过期时间；
- WithCancel 函数，是将父节点复制到子节点，并且返回一个额外的 CancelFunc 函数类型变量。在父 goroutine 中，通过 WithCancel 可以创建子节点的 Context, 还获得了子 goroutine 的控制权，一旦执行了 CancelFunc 函数，子节点 Context 就结束了；
- WithDeadline 函数，也是将父节点复制到子节点，但是其过期时间是由 deadline 和 parent 的过期时间共同决定。当 parent 的过期时间早于 deadline 时，返回的过期时间与 parent 的过期时间相同。父节点过期时，所有的子孙节点必须同时关闭；
- WithTimeout 函数，传入的是从现在开始 Context 剩余的生命时长。也都返回了所创建的子 Context 的控制权，一个 CancelFunc 类型的函数变量；
- WithValue 函数，返回 parent 的一个副本，调用该副本的 Value(key) 方法将得到 value；

原则：
1. 不要把 context 放到一个结构体中，应该作为第一个参数显式地传入函数；
2. 即使方法允许，也不要传入一个 nil 的 context，如果不确定需要什么 context 的时候，传入一个 context.TODO；
3. 使用 context 的 Value 相关方法应该传递和请求相关的元数据，不要用它来传递一些可选参数；
4. 同样的 context 可以传递到多个 goroutine 中，Context 在多个 goroutine 中是安全的；
5. 在子 context 传入 goroutine 中后，应该在子 goroutine 中对该子 context 的 Done channel 进行监控，一旦该 channel 被关闭，应立即终止对当前请求的处理，并释放资源。

## channel 相关

### channel 是否线程安全，锁用在什么地方

线程安全：不同协程通过 channel 进行通信，本身使用场景就是多线程，为了保证数据的一致性，必须实现线程安全。
channel 的底层实现中，hchan 结构体中采用 Mutex 锁来保证数据读写安全。在对循环数组 buf 中的数据进行入队和出队操作时，必须先获取互斥锁，才能操作 channel 数据。

### channel 的底层实现原理（数据结构）

hchan: 循环数组 buf, 下一个要发送数据的下标 sendx, 下一个要接收数据的下标 recvx, 发送队列 sendq, 接收队列 recvq, 互斥锁 lock。

### nil、关闭的 channel、有数据的 channel，再进行读、写、关闭会怎么样

两种类型：无缓冲、有缓存

三种模式：
1. 写操作模式（单向通道）：make(chan<- int)
2. 读操作模式（单向通道）：make(<-chan int)
3. 读写操作模式（双向通道）：make(cha int)

三种状态：
            关闭     发送            接收
1. 未初始化： panic   永远阻塞导致死锁  永远阻塞导致死锁
2. 正常：    正常关闭  阻塞或者成功发送  阻塞或者成功接收
3. 关闭：    panic    panic          缓冲区为空则为零值，否则可以继续读

注意：
1. 一个 channel 不能多次关闭，会导致 panic；
2. 若多个 goroutine 监听同一个 channel，那么 channel 上的数据可能随机被某一个 goroutine 取走消费；
3. 若多个 goroutine 监听同一个 channel，如果这个 channel 被关闭，则所有 goroutine 都能接收到退出信号。

非阻塞队列
```go
func push(q chan int, item int) error {
	select {
	case q <- item:
		return nil
	default:
		return errors.New("queue full")
	}
}

func get(q chan int) (int, error) {
	var item int
	select {
	case item = <-q:
		return item, nil
	default:
		return 0, errors.New("queue empty")
	}
}

func TestNonBlockingQueue(t *testing.T) {
	q := make(chan int, 5)
	x := []int{1, 2, 3, 4, 5, 6}
	for _, value := range x {
		err := push(q, value)
		fmt.Printf("error:%v\n", err)
	}

	for _, value := range x {
		fmt.Println(value)
		v, err := get(q)
		fmt.Printf("v:%v, error:%v\n", v, err)
	}
}
```

带超时的阻塞队列
```go
func push(q chan int, item int, timeoutSecs int) error {
	select {
	case q <- item:
		return nil
	case <-time.After(time.Duration(timeoutSecs) * time.Second):
		return errors.New("queue full, wait timeout")
	}
}

func get(q chan int, timeoutSecs int) (int, error) {
	var item int
	select {
	case item = <-q:
		return item, nil
	case <-time.After(time.Duration(timeoutSecs) * time.Second):
		return 0, errors.New("queue empty, wait timeout")
	}
}

func TestTimeoutBlockingQueue(t *testing.T) {
	q := make(chan int, 5)
	x := []int{1, 2, 3, 4, 5, 6}
	for _, value := range x {
		err := push(q, value, 3)
		fmt.Printf("error:%v\n", err)
	}

	for _, value := range x {
		fmt.Println(value)
		v, err := get(q, 3)
		fmt.Printf("v:%v, error:%v\n", v, err)
	}
}
```

### 向 channel 发送数据和从 channel 读取数据的流程

发送流程：

阻塞式：调用 chansend 函数，并且 block=true
```go
ch <- 10
```

非阻塞式：调用 chansend 函数，并且 block=false
```go
select {
case ch <- 10:
	...
default:
	...
}
```

向 channel 中发送数据时分为两大块：检查和数据发送，流程如下：
- 若 channel 的读等待队列存在接收者 goroutine
  - 将数据直接发送给第一个等待的 goroutine，唤醒接收的 goroutine
- 若 channel 的读等待队列不存在接收者 goroutine
  - 若循环数组 buf 未满，那么将数据发送到 buf 的队尾
  - 若循环数组 buf 已满，就走阻塞发送流程，将当前 goroutine 加入写等待队列，并挂起等待唤醒

接收流程：

向 channel 中接收数据时分为两大块：检查和数据接收，流程如下：
- 若 channel 的写等待队列存在发送者 goroutine
  - 若是无缓冲 channel，直接从第一个发送者 goroutine 把数据拷贝到接收变量，唤醒发送的 goroutine
  - 若是有缓冲 channel (已满)，将循环数组 buf 的队首元素拷贝给接收变量，将第一个发送者 goroutine 的数据拷贝到 buf 队尾，唤醒发送的 goroutine
- 若 channel 的写等待队列不存在发送者 goroutine
  - 若循环数组非空，将 buf 的队首元素拷贝给接收变量
  - 若循环数组为空，就走阻塞接收流程，将当前 goroutine 加入读等待队列，并挂起等待唤醒

### channel 底层数据结构和主要使用场景

```go
type hchan struct {
	qcount   uint // 数组长度
	dataqsiz uint // 数组容量
	buf      unsafe.Pointer // 数组地址
	elemsize uint16 // 元素大小
	closed   uint32 // 关闭状态
	elemtype *_type // 元素类型
	sendx    uint // 下一次写下标位置
	recvx    uint // 下一次读下标位置
	sendq    waitq // 写等待队列
	recvq    waitq // 读等待队列
	lock     mutex // 互斥锁，不允许并发读写
}
```

无缓冲和有缓冲区别：
- 管道没有缓冲区，从管道读数据会阻塞，直到有协程向管道中写入数据。同样，向管道写入数据也会阻塞，直到有协程从管道读取数据；
- 管道有缓冲区但缓冲区没有数据，从管道读取数据也会阻塞，直到协程写入数据，如果管道满了，写数据也会阻塞，直到协程从缓冲区读取数据。

特点：
1. 读写值 nil 管道会永久阻塞；
2. 关闭的管道仍可以读数据；
3. 往关闭的管道写数据会 panic；
4. 关闭为 nil 的管道 panic；
5. 关闭已经关闭的管道 panic。

使用场景：消息传递、消息过滤，信号广播，事件订阅与广播，请求、响应转发，任务分发，结果汇总，并发控制，限流，同步与异步。

## GMP 相关

### 什么是 GMP

G: goroutine 协程
M: thread 线程
P: processor 上下文处理器

golang 中线程是运行 goroutine 的实体，调度器的作用是把可运行的 goroutine 分配到工作线程上。

- 全局队列：存放等待运行的 G；
- P 的本地队列：也存放等待运行的 G，存的数量不超过 256 个。新建 G' 时，G' 优先加入 P 的本地队列，若队列满了，则把本地队列一半的 G 移到全局队列；
- P 列表：所有的 P 都在程序启动时创建，并保存在数组中，最多有 GOMAXPROCS (可配置) 个；
- M：线程想运行任务就得获取 P，从 P 的本地队列获取 G，P 队列为空时，M 也会尝试从全局队列拿一批 G 放到 P 的本地队列，或从其他 P 的本地队列偷一半放到自己 P 的本地队列。M 运行 G，G 执行之后，M 会从 P 获取到下一个 G，不断重复下去。

Goroutine 的调度器和 OS 调度器是通过 M 结合起来的，每个 M 都代表了 1 个内核线程，OS 调度器负责把内核线程分配到 CPU 的核上执行。

- 可以通过 go func () 创建一个 goroutine；
- 有两个存储 G 的队列，一个是调度器 P 的本地 G 队列、一个是全局 G 队列。新创建的 G 会先保存在 P 的本地队列，如果 P 的本地队列已满就会保存在全局的队列里；
- G 只能运行在 M 中，一个 M 必须持有一个 P，M 与 P 是 1：1 的关系。M 会从 P 的本地队列弹出一个可执行状态的 G 来执行，如果 P 的本地队列为空，就会从其他 MP 组合偷取一个可执行的 G 来执行；
- 一个 M 调度 G 执行的过程是一个循环机制；
- 当 M 执行某一个 G 时候如果发生了 syscall（系统调用） 等操作，M 会阻塞，如果当前正好有一些 G 在执行，runtime 会把这个线程 M 从 P 中摘除，然后再创建一个新的操作系统的线程 (如果有空闲的线程可用就复用空闲线程) 来服务于这个 P；
- 当 M 系统调用结束时，这个 G 会尝试获取一个空闲的 P 执行，并放入到这个 P 的本地队列。如果获取不到 P，那么这个线程 M 变成休眠状态， 加入到空闲线程中，然后这个 G 会被放入全局队列中。

关于 G,P,M 的个数问题，G 的个数理论上是无限制的，但是受内存限制，P 的数量一般建议是逻辑 CPU 数量的 2 倍，M 的数据默认启动的时候是 10000，内核很难支持这么多线程数，所以整个限制客户忽略，M 一般不做设置，设置好 P，M 一般都是要大于 P。

### 进程、线程、协程的区别

- 进程：是应用程序的启动实例，每个进程都有独立的内存空间，不同的进程通过进程间的通信方式来通信；
- 线程：每个进程至少包含一个线程，是 CPU 调度的基本单位，多个线程之间可以共享进程的资源并通过共享内存等线程间的通信方式来通信；
- 协程：轻量级线程，不受操作系统的调度，协程的调度器由用户应用程序提供，协程调度器按照调度策略把协程调度到线程中运行。

### 抢占式调度是如何抢占的

- 基于协作的抢占式调度
- 基于信号量的抢占式调度

## 锁相关

### 除了 mutex 还有哪些方式安全读写共享变量

- 将共享变量放在一个 goroutine 中，其他 goroutine 通过 channel 进行读写操作；
- 可以用个数为 1 的信号量 semaphore 实现互斥；

### Go 如何实现原子操作

Go 的标准库代码包 sync/atomic 提供了原子的读取 (Load 为前缀的函数) 或写入 (Restore 为前缀的函数)。

原子操作与互斥锁的区别：
- 互斥锁是一种数据结构，用来让一个线程执行程序的关键部分，完成互斥的多个操作；
- 原子操作是针对某个值的单个互斥操作。

### Mutex 是悲观锁还是乐观锁

Mutex 是悲观锁，互斥锁。

锁的实现一般依赖于信号量，是一个非负的整数计数器。

- 信号量：多线程同步使用的；一个线程完成某个动作后通过信号告诉别的线程，别的线程才可以执行某些动作；非负整数；
- 互斥量：多线程互斥使用的；一个线程占用某个资源，别的线程无法访问，直至该线程离开；0 或 1

- 悲观锁：互斥锁。借助数据库锁机制，在修改数据之前先锁定，再修改的方式称为悲观并发控制 Pessimistic Concurrency Control PCC
  - 加锁，就是把信号量减 1，若是 0 则加锁成功；释放锁时把信号量加 1，若是 1 则释放成功。
- 乐观锁：读写锁。假定数据一般情况下不会造成冲突，在数据进行提交更新时，才会真正对数据的冲突与否进行检测，若冲突则返回异常信息，让用户决定如何去做。适用于读多写少的场景。

### Mutex 有几种模式

- 正常模式：
  - 当前 mutex 只有一个 goroutine 来获取，没有竞争，直接返回；
  - 新的 goroutine 进来，若当前 mutex 已被获取，则该 goroutine 进入一个先进先出的 waiter 队列，在 mutex 被释放后，waiter 按照先进先出的方式获取锁。该 goroutine 会处于自旋状态。(不挂起，继续占用 CPU)；
  - 新的 goroutine 进来，mutex 处于空闲状态，将参与竞争。新来的 goroutine 有先天的优势，它们正在 CPU 中运行，可能它们的数量还不少，所以，在高并发情况下，被唤醒的 waiter 可能比较悲剧地获取不到锁，这时，它会被插入到队列的前面。如果 waiter 获取不到锁的时间超过阈值 1 毫秒，那么，这个 Mutex 就进入到了饥饿模式。
- 饥饿模式：
  - 在饥饿模式下，Mutex 的拥有者将直接把锁交给队列最前面的 waiter。新来的 goroutine 不会尝试获取锁，即使看起来锁没有被持有，它也不会去抢，也不会 spin（自旋），它会乖乖地加入到等待队列的尾部。 如果拥有 Mutex 的 waiter 发现下面两种情况的其中之一，它就会把这个 Mutex 转换成正常模式:
    - 此 waiter 已经是队列中的最后一个 waiter 了，没有其它的等待锁的 goroutine 了；
    - 此 waiter 的等待时间小于 1 毫秒。

### goroutine 的自旋占用资源如何解决

自旋锁：当一个线程在获取锁的时候，若锁已被其他线程获取，那么该线程会循环等待，不断判断锁是否能被成功获取，直到获取到锁才会退出循环。

自旋的条件：
- 还没自旋超过 4 次；
- 多核处理器；
- GOMAXPROCS > 1;
- P 上本地 goroutine 队列为空。

mutex 会让当前的 goroutine 去空转 CPU，在空转完后再次调用 CAS 方法去尝试性的占有锁资源，直到不满足自旋条件，则最终会加入到等待队列里。

## 并发相关

### 如何控制并发数

一、有缓冲通道

根据通道中没有数据时读取操作陷入阻塞和通道已满时继续写入操作陷入阻塞的特性，正好实现控制并发数量。

```go

```




























