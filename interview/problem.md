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

