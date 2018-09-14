package main

import (
     "astuart.co/go-sse"
    "fmt"
    "io/ioutil"
    "time"
    "log"
)
func main() {
    fmt.Println("Hello")

    uri:="http://localhost:3000/"
    ch := make(chan *sse.Event , 100)
    go sse.Notify(uri, ch)

    fmt.Println("Start main loop")
    for {
        select {
        case msg := <-ch:
            fmt.Println("received message")

            b, err := ioutil.ReadAll(msg.Data)
            if err != nil {
                log.Fatal(err)
	        }
            fmt.Printf("%s", b)
    default:
        fmt.Println("no message received")
    }

       time.Sleep(2 * time.Second)
    }



}
