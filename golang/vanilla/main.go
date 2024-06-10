package main

import (
    "flag"
    "fmt"
//    "crypto/md5"
    "net/http"
//    "os"
    "runtime"
//    "strconv"
//    "syscall"
//    "time"
//    "encoding/json"

    "github.com/devhands-io/bootcamp-samples/golang/vanilla/handlers"
)

var (
    data []byte
)

var (
    host string
    port int
)

func init() {
    flag.StringVar(&host, "host", "localhost", "server host")
    flag.IntVar(&port, "port", 8000, "server port")
}

func main() {
    runtime.GOMAXPROCS(2 * runtime.NumCPU())

    flag.Parse()

    http.HandleFunc("/", handlers.Ok)
    http.HandleFunc("/hello", handlers.Hello)

/*
    // sleeps
    cpuSleep := NewGetrusagePayload()
    ioSleep := NewIOPayload()
    http.HandleFunc("/payload", SleepHandler(cpuSleep, ioSleep))
*/
    addr := fmt.Sprintf("%s:%d", host, port)
    fmt.Println("serving at " + addr)
    http.ListenAndServe(addr, nil)
}
