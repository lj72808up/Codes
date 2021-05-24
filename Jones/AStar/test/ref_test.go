package test

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestRef(t *testing.T) {
	var v AInter
	v = AStru{}
	v.SayHello() // 可以
	//v.SayNo()     // 不可以, v没有SayNo方法
}

type AInter interface {
	SayHello() string
}

type AStru struct{}

func (this AStru) SayHello() string {
	return "my name is Astru"
}

func (this AStru) SayNo() string {
	return "No!"
}

func Test2(t *testing.T) {
	var r io.Reader
	file, _ := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	r = file

	var w io.Writer
	w = r.(io.Writer)

	fmt.Println(w)
}

func Test3(t *testing.T) {
	var x float64 = 3.4
	fmt.Println("type:", reflect.TypeOf(x))            // type: float64
	fmt.Println("value:", reflect.ValueOf(x).String()) // value: <float64 Value>
	fmt.Println("value:", reflect.ValueOf(x).Float())  // value: 3.4

	v := reflect.ValueOf(x)
	y := v.Interface().(float64)
	fmt.Println(y)
	fmt.Printf("value is %7.1e\n", v.Interface())
}

func Test4(t *testing.T) {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("settability of v:", v.CanSet()) // settability of v: false

	var x2 float64 = 3.4
	v2 := reflect.ValueOf(&x2)
	fmt.Println("type of v2:", v2.Type())                       // type of v2: *float64
	fmt.Println(fmt.Println("settability of v2:", v2.CanSet())) // settability of v2: false

	e2 := v2.Elem()
	fmt.Println("type of e2:", e2.Type())                       // type of e2: float64
	fmt.Println(fmt.Println("settability of e2:", e2.CanSet())) // settability of e2: true
	e2.SetFloat(4.1)
	fmt.Println(e2.Interface()) // 4.1
	fmt.Println(x2)             // 4.1
}

func Test5(t *testing.T) {
	tt := TT{23, "skidoo"}
	s := reflect.ValueOf(&tt).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	// output:
	// 0: A int = 23
	// 1: B string = skidoo

	s.Field(0).SetInt(77)
	s.Field(1).SetString("Sunset Strip")
	fmt.Println("tt is now", tt)   // tt is now {77 Sunset Strip}
}

type TT struct {
	A int
	B string
}
