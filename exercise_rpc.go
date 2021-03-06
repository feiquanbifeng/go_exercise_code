// rpc call
// use default encode/decode gob

package main

import (
    "errors"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/rpc"
)

type Def struct {
    Name string
}

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}

func main() {
    arith := new(Arith)
    rpc.Register(arith)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":1234")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    go http.Serve(l, nil)

    client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }
    args := &Args{7, 8}

    // Synchronous call
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
    fmt.Println()

    // Asynchronous call
    quotient := new(Quotient)
    divCall := client.Go("Arith.Divide", args, quotient, nil)
    replyCall := <-divCall.Done

    if replyCall.Error != nil {
        log.Fatal("divide error:", replyCall)
    }
    fmt.Printf("Quotient:%d/%d=%d, %d%%%d=%d", args.A, args.B, quotient.Quo, args.A, args.B, quotient.Rem)
}
