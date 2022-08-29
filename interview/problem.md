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

### rune

- rune: uint32，处理 unicode 或 utf-8 字符，区分字符值和整数值。
- byte: uint8，处理 ascii 字符。

### interface

- 它是一个方法的集合，但没有方法的实现，也没有数据字段；
- 在函数内部，只在乎传入的实参有没有实现形参接口的全部方法，实现了就能传入；
- go 是静态语言，在编译阶段就能检测出赋值给接口的值，有没有实现该接口的所有方法；python 是动态语言，需要运行期间才能检测出来。

#### 值接收者和指针接收者（值调用者和指针调用者）

- 给用户自定义类型添加新的方法时，与函数的区别是，需要给函数前添加一个接收者，既可以是自定义类型的值类型，也可以是自定义类型的指针类型；
- 在调用方法时，`值类型`既可以调用`值接收者`的方法，也可以调用`指针接收者`的方法；`指针类型`既可以调用`指针接收者`的方法，也可以调用`值接收者`的方法；
- 在方法内部，如果对接收者进行修改，无论是值类型调用还是指针类型调用，只有当接收者是`指针类型`时，才会影响到接收者；
- 对自定义类型实现接口的方法时，注意：
  - 实现接口方法的接收者/
  - 赋值给接口的类型      值接收者   指针接收者
  - 接口=值类型          可以      报错
  - 接口=指针类型        可以      可以

*类型不同可以调用*

```go
type field struct {
    name string
}

func (p *field) pointerMethod() {
    fmt.Println(p.name)
}

func (p field) valueMethod() {
    fmt.Println(p.name)
}

func main() {
    fp := &field{name: "pointer"}
    fv := field{name: "value"}
    
    fp.pointerMethod()
    fp.valueMethod()
    fv.pointerMethod()
    fv.valueMethod()
}
```
- 在值类型调用指针接收者方法时，实际为 `(&fv).pointerMethod()`;
- 在指针类型调用值接收者方法时，实际为 `(*fp).valueMethod()`.

*类型不同不可以调用*

有两种情况：都是`值类型`不能调用`指针接收者方法`。
- `值类型`不能被寻址；
- 用`指针接收者`实现接口。

`值类型`不能被寻址：
```go
type field struct {
	name string
}

func (p *field) pointerMethod() {
	fmt.Println(p.name)
}

func (p field) valueMethod() {
	fmt.Println(p.name)
}

func NewField() field {
	return field{name: "right value struct"}
}

func main() {
	NewField().valueMethod()
	NewField().pointerMethod()
}
```
编译器首先给 NewField() 返回的右值调用 pointer method，出错；然后试图给其插入取地址符，报错。

区别是是否可以被寻址：
- 左值：可以被寻址，既可以出现在赋值号左边也可以出现在右边；
- 右值：不可以被寻址，如函数返回值、字面值、常量值等，只能出现在赋值号右边。

用`指针接收者`实现接口：
```go
type human interface {
	speak()
	sing()
}

type man struct{}

func (m man) speak() {
	fmt.Println("speaking")
}

func (m *man) sing() {
	fmt.Println("singing")
}

func main() {
	var h human = &man{}
	//var h human = man{}
	h.speak()
	h.sing()
}
```
- 如果是值接收者，实体类型的值和指针都可以实现对应的接口；如果是指针接收者，只有类型的指针能够实现接口；
- 接收者是指针类型的方法，很可能在方法中对接收者的属性进行修改；但接收者是值类型的方法，在方法中无法影响接收者本身；
- 如果实现了接收者是值类型的方法，会隐含地实现了接收者是指针类型的方法。

#### 接口的类型检查

*断言*

- <目标类型的值>, <布尔参数> := <表达式>.(目标类型) // 安全类型断言
- <目标类型的值> := <表达式>.(目标类型) // 非安全类型断言

```go
type Student struct {
	Name string
	Age  int
}

func main() {
	stu := &Student{
		Name: "小有",
		Age:  22,
	}

	var i interface{} = stu
	s1 := i.(*Student) // 断言成功，s1 为 *Student 类型，不安全断言
	fmt.Println(s1)

	s2, ok := i.(Student) // 断言失败，ok 为 false，安全断言
	if ok {
		fmt.Println("success:",s2)
	}
	fmt.Println("failed:",s2)
}
```

*接口类型有多种情况，采用 Type Switch 方法*

```go
func typeCheck(v interface{}){
	//switch v.(type) {       // 只用判断类型，不需要值
	switch msg := v.(type) {   // 值和判断类型都需要
		case int :
			...
		case string:
			...
		case Student:
			...
		case *Student:
			...
		default:
			...
	}
}
```

`fmt.Println` 参数是 `interface`，其打印机制是什么
- 若为内置类型，则会穷举真实类型，然后打印；
- 若为自定义类型，会先检查是否实现 String() 方法，若实现了直接调用，否则利用反射来遍历对象成员，进行打印。
  - 注：别在自定义类型的 String() 方法中 fmt.Println 自己，会造成递归打印。

```go
type base interface{ F() }

type student struct{ Name string }

func (s *student) F() {}

type class struct{ Name string }

func (c *class) F() {}

type teacher struct{ Name string }

func (t *teacher) F() {}

func isType(v interface{}) {
	switch msg := v.(type) {
	case student, teacher:
		fmt.Println(msg.Name) // 这里会报错，因为 msg 是 interface 类型，没有 Name 属性
	case class:
		fmt.Println(msg.Name) // 这里不会报错，因为 msg 是 class 类型，有 Name 属性
	}
}
```
- switch type 的 case 后面只有一个类型 T1，则 msg 对应的类型是 T1；
- switch type 的 case 后面有多个类型 (T2,T3)，则 msg 对应的类型是 interface.

#### 空接口

- 空接口 interface{} 不包含任何的 method，故所有类型都实现了空接口；
- 空接口对于描述起不到任何作用，但在需要存储任意类型的数值时非常有用，它可以存储任意类型的数值。

```go
type Student struct {
}

func Set(x interface{}) {
}

func Get(x *interface{}) {
}

func main() {
 s := Student{}
 p := &s
 // A B C D
 Set(s)
 Get(s)
 Set(p)
 Get(p)
}
```
- B、D 会报错。只能放入接口的指针类型 *interface{}.

### 反射

#### Go 的反射包

```go
ty := reflect.TypeOf(Person) // 获取类型
mt1 := ty.Method(0) // 获取第几个方法
mt2, _ := ty.MethodByName("Get") // 根据方法名称获取方法
```

#### DeepEqual 的作用和原理

不能用 == 比较的情况：
- 切片、map、函数
- 含有以上三种的结构体和数组

声明两个比较值的结构体的名称不同，即使字段名、类型、顺序相同，也不能比较（强转类型可以比较），必须用同一个结构体声明的值，才能比较。

作用：判断两个变量的实际内容完全一致。
- Array：相同索引处的元素"深度"相等；
- Struct：相应字段，包括导出和不导出，"深度"相等；
- Func：只有两者都是 nil；
- Interface：两者存储的具体值"深度"相等；
- Map：1、都为 nil；2、非空，长度相等，指向同一个 map 实体对象，或者相应的 key 指向的 value "深度"相等；
- Pointer：1、使用 == 比较的结果为真；2、指向的实体"深度"相等；
- Slice：：1、都为 nil；2、非空，长度相等，首元素指向同一个底层数组的相同元素，即 &x[0] == &y[0]，或者相同索引处的元素"深度"相等；
- numbers, bools, strings, channels：使用 == 比较的结果为真。

原理：
- 先会判断两者中是否存在 nil，当两者都是 nil，返回 true；
- 利用反射，获取两者类型，若不相同返回 false；
- 调用 deepValueEqual 函数判断，它起始是一个递归函数，一直递归到最基本的类型使用 == 去判断，然后层层返回，得出比较结果；
- 对于特殊情况，如 map、slice，会先判断长度是否相等（不同返回 false），指针是否相等（相同返回 true），然后再判断里面的具体值是否相同，其实是一个快速对比处理。

### init 函数

在包初始化的时候会调用 init 函数，不能被显式调用。

适用场景：
- 初始化变量；
- 检查或修复程序状态；
- 注册任务；
- 仅需要执行一次的情况。

特征：
- 同一个包，可以有多个 init 函数；
- 包中的每个源文件中可以有多个 init 函数。

执行顺序：
- 同一个源文件的 init 函数，是按照先后顺序执行（且在全局变量初始化之后）；
- 同一个包中的源文件的 init 函数，是按照源文件的字母顺序执行；
- 不同包的 init 函数，按照包导入的依赖关系决定先后顺序。

### sync.WaitGroup 函数

一个 WaitGroup 对象，可以实现同一时间启动 n 个协程，并发执行，等 n 个协程全部执行结束后，再继续往下执行的功能。
通过 Add() 方法设置启动多少个协程，在每个协程结束后执行 Done() 方法，计数减一，同时用 Wait() 方法阻塞主协程，等待所有的协程执行结束。

- Add 方法和 wait 方法不可以并发同时调用，Add 方法要在 wait 方法之前调用；
- Add() 设置的值必须与实际等待的 goroutine 数量一致，否则会 panic；
- 调用了 wait 方法后，必须在 wait 返回之后才能再重新使用 waitGroup，也就是 wait 返回之前不要调用 Add 方法，否则会 panic；
- Done 只是对 Add 的简单封装，可以向 Add 传入任意负数，快速将计数器归零，以唤醒等待的 goroutine；
- waitGroup 对象只能有一份，不可以拷贝给其他变量，否则会造成 bug；
- waitGroup 对象不是一个引用类型，通过函数传值的时候需要使用地址，因为 Go 只有值传递。

Go 语言中提供了两种 copy 检查，一种是在运行时进行检查，一种是通过静态检查。
不过运行检查是比较影响程序的执行性能的，Go 官方目前只提供了 strings.Builder 和 sync.Cond 的 runtime 拷贝检查机制，对于其他需要 nocopy 对象类型来说，使用 go vet 工具来做静态编译检查。运行检查的实现可以通过比较所属对象是否发生变更。

## 切片篇

- nil 切片：
```go
var slice []int
slice := *new([]int)
```

- 空切片：
```go
slice := []int{}
slice := make([]int, 0)
```



































