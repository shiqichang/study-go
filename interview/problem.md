## 特性篇

### Golang 使用什么数据类型

布尔型 bool、数值型（整型 int、浮点型 float）、字符串 string
指针 pointer、数组 [...]int、切片 []int、结构体 struct、映射 map、管道 chan、接口 interface、函数 func

### 数组定义问题

数组是可以指定下标的方式定义。

```go
array := [...]int{1,2,3,9:34} // array[9]=34, len(array)=10

m := [...]int{'a':1,'b':2,'c':3} // len(m)=100，因为 'c' 的 ascii 是 99
```

### 内存四区

- 代码区：存放代码；
- 全局区：常量+全局变量。进程退出时由操作系统回收；
- 堆区：空间充裕，数据存放时间较久。一般由开发者分配，启动 Golang 的 GC 由 GC 清除机制自动回收；
- 栈区：空间较小，要求数据读写性能高，数据存放时间较短。由编译器自动分配和释放，存放函数的参数值、局部变量、返回值等。

### 空结构体

空结构体：不包含任何字段的结构体 struct{}

```go
var et struct{}
et := struct{}{}
type ets struct{}
et := ets{}
var et ets
```

- 所有空结构体的地址都是同一地址，zerobase 地址，大小为 0；
- 用于保存不重复的元素的集合，如 map，struct{} 作为 value 不占用额外空间；
- 用于 channel 中信号传输，不在乎传输内容，只要信号；
- 作为方法的接收者，空结构体嵌到其他结构体中，实现继承。

### 如何停止一个 goroutine

- for-select 方法，采用通道，通知协程退出；
- 采用 context 包。

### cap 函数可以作用于什么

数组、切片、通道。

### Printf()、Sprintf()、Fprintf() 的区别

- Printf()：标准输出，用于打印；
- Sprintf()：把格式化字符串输出到字符串，并返回；
- Fprintf()：把格式化字符串输出到实现了 io.Writer() 方法的类型，比如文件。

### Go 中值传递和引用传递

Go 中的函数传参都是值传递：值的副本或指针的副本。

int：值传递
```go
func TestInt(t *testing.T) {
	a := 100
	var b, c = &a, &a
	fmt.Println(b, c)   // 0x14000717968 0x14000717968
	fmt.Println(&b, &c) // 0x140003020b0 0x140003020b8

	d := 200
	b = &d
	fmt.Println(a, *b, *c) // 100 200 100
	// b 和 c 都保存了 a 的地址，但 b、c 本身是独立的，改变 b 的值不会对 a、c 产生影响
}
```

`go tool compile -S -N -l xxx.go` 打印汇编信息。

### 常量计数器 iota

用常量定义代替枚举类型。

```go
const (
    mutexLocked = 1 << iota  // 1
    mutexWoken               // 2
    mutexStarving            // 4
    mutexWaiterShift = iota  // 3
)
```

- 不同 const 定义块互不干扰；
- 所有注释行和空白行全部忽略；_ 代表一行，不能忽略；
- 没有表达式的常量定义复用上一行的表达式；
- 从第一行开始，iota 从 0 逐行加一；这一行即使没有iota也算一行；
- 替换所有 iota。

### defer

```go
package main
import "fmt"
func returnButDefer() (t int) {  // t 初始化 0，并且作用域为该函数全域
    defer func() {
        t = t * 10
    }()
    return 1
}
func main() {
    fmt.Println(returnButDefer())   // 输出 10
}
```

```go
func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc2(i int) int {
	t := i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc3(i int) (t int) {
	defer func() {
		t += i
	}()
	return 2
}

func DeferFunc4() (t int) {
	defer func(i int) {
		fmt.Println(i)
		fmt.Println(t)
	}(t)
	t = 1
	return 2
}

func main() {
	fmt.Println(DeferFunc1(1))
	fmt.Println(DeferFunc2(1))
	fmt.Println(DeferFunc3(1))
	DeferFunc4()
}
```


















