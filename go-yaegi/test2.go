package main

import (
	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"reflect"
)

const payload = `
package main

import (
	"fmt"
	"github.com/tidwall/pretty"

	"log"
	"os/exec"
)

func main(){
	var example2 = "{\"name\": \"hello\"}"
	result := pretty.Pretty([]byte(example2))
	fmt.Println(string(result))
}
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

	i.Use(stdlib.Symbols)


	EvalWithCheck(i, payload)

}

