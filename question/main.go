package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
)

const (
	Left = iota
	Top
	Right
	Bottom
)

func main() {
	//alternatePrint()

	//fmt.Println(isUniqueString("string"))
	//fmt.Println(isUniqueString2("strings"))

	//fmt.Println(revertString("abcdefg"))

	//fmt.Println(isRegroup("abc", "bac"))

	//fmt.Println(replaceBlank("ab c"))

	//fmt.Println(move("R2(LF)", 0, 0, Top))

	//q007()

	q008()
}

func alternatePrint() {
	letter, number := make(chan bool), make(chan bool)
	wait := sync.WaitGroup{}

	go func() {
		i := 1
		for {
			<- number
			fmt.Printf("%d%d", i, i + 1)
			i += 2
			letter <- true
		}
	}()
	wait.Add(1)
	go func(wait *sync.WaitGroup) {
		i := 'A'
		for {
			<- letter
			if i >= 'Z' {
				wait.Done()
				return
			}
			fmt.Printf("%s%s", string(i), string(i + 1))
			i += 2
			number <- true
		}
	}(&wait)
	number <- true
	wait.Wait()
}

func isUniqueString(s string) bool {
	if strings.Count(s, "") > 3000 {
		return false
	}

	for _, v := range s {
		if v > 127 {
			return false
		}
		if strings.Count(s, string(v)) > 1 {
			return false
		}
	}

	return true
}

func isUniqueString2(s string) bool {
	if strings.Count(s, "") > 3000 {
		return false
	}

	for k, v := range s {
		if v > 127 {
			return false
		}
		if strings.Index(s, string(v)) != k {
			return false
		}
	}

	return true
}

func revertString(s string) (string, bool) {
	str := []rune(s)
	l := len(str)
	if l > 5000 {
		return s, false
	}

	for i := 0; i < l/2; i++ {
		str[i], str[l-1-i] = str[l-1-i], str[i]
	}
	return string(str), true
}

func isRegroup(s1, s2 string) bool {
	sl1 := len([]rune(s1))
	sl2 := len([]rune(s2))

	if sl1 > 5000 || sl2 > 5000 || sl1 != sl2 {
		return false
	}

	for _, v := range s1 {
		if strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
			return false
		}
	}

	return true
}

func replaceBlank(s string) (string, bool) {
	if len([]rune(s)) > 1000 {
		return s, false
	}

	for _, v := range s {
		if string(v) != " " && unicode.IsLetter(v) == false {
			return s, false
		}
	}
	return strings.Replace(s, " ", "%20", -1), true
}

func move(cmd string, x0, y0, z0 int) (x, y, z int) {
	x, y, z = x0, y0, z0
	repeat := 0
	repeatCmd := ""

	for _, s := range cmd {
		switch {
		case unicode.IsNumber(s):
			repeat = repeat * 10 + (int(s) - '0')
		case s == ')':
			for i := 0; i < repeat; i++ {
				x, y, z = move(repeatCmd, x, y, z)
			}
			repeat = 0
			repeatCmd = ""
		case repeat > 0 && s != '(' && s != ')':
			repeatCmd += string(s)
		case s == 'L':
			z = (z + 1) % 4
		case s == 'R':
			z = (z - 1 + 4) % 4
		case s == 'F':
			switch {
			case z == Left || z == Right:
				x = x - z + 1
			case z == Top || z == Bottom:
				y = y - z + 2
			}
		case s == 'B':
			switch {
			case z == Left || z == Right:
				x = x + z - 1
			case z == Top || z == Bottom:
				y = y + z - 2
			}
		}
	}
	return
}

func q007() {
	type Param map[string]interface{}
	type Show struct {
		Param
	}
	s := new(Show)
	s.Param = Param{}
	s.Param["RMB"] = 10000
	fmt.Println(s)

	type People struct {
		name string `json:"name"`
	}
	js := `{
        "name":"11"
    }`
	var p People
	err := json.Unmarshal([]byte(js), &p)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("people: ", p)

	var value int32
	v := value
	if atomic.CompareAndSwapInt32(&value, v, v+1) {
	}

	type Student struct {
		name string
	}
	m := map[string]*Student{"people": {"zhoujielun"}}
	fmt.Println(m["people"].name)
	m["people"].name = "wuyanzu"
	fmt.Println(m["people"].name)

	ret := exec("111", func(n string) string {
		return n + "func1"
	}, func(n string) string {
		return n + "func2"
	}, func(n string) string {
		return n + "func3"
	}, func (n string) string {
		return n + "func4"
	})
	fmt.Println(ret)
}

type query func(string) string

func exec(name string, vs ...query) string {
	ch := make(chan string)
	fn := func(i int) {
		ch <- vs[i](name)
	}
	for i, _ := range vs {
		go fn(i)
	}
	return <-ch
}

func q008() {
	deferCall()
}

func deferCall() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}
