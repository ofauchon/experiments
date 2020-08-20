package main

import (
	"github.com/containous/yaegi/interp"
	"reflect"
)

const src = `package foo
func Bar1(s string) string { return s + "-Foo1" }
func Bar2(s string) string { return s + "-Foo2" }
func Bar3(s string) string { return s + "-Foo3" }
`

func EvalWithCheck(i *interp.Interpreter, src string) reflect.Value {
	v, err := i.Eval(src)
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	i := interp.New(interp.Options{})

	EvalWithCheck(i, src)

	bar1 := EvalWithCheck(i, "foo.Bar1").Interface().(func(string) string)
	bar2 := EvalWithCheck(i, "foo.Bar2").Interface().(func(string) string)
	bar3 := EvalWithCheck(i, "foo.Bar3").Interface().(func(string) string)

	r1 := bar1("Kung")
	r2 := bar2("Kung")
	r3 := bar3("Kung")
	println(r1,r2,r3)
}

