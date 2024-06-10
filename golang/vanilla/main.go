package main

import (
    "flag"
    "fmt"
    "net/http"
    "runtime"
    "github.com/devhands-io/bootcamp-samples/golang/vanilla/handlers"
    "github.com/devhands-io/bootcamp-samples/golang/vanilla/payload"
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

    // dummy handlers
    http.HandleFunc("/", handlers.Ok)
    http.HandleFunc("/hello", handlers.Hello)

    // payload
    cpuSleep := payload.NewGetrusagePayload()
    ioSleep := payload.NewIOPayload()
    http.HandleFunc("/payload", handlers.SleepHandler(cpuSleep, ioSleep))

    addr := fmt.Sprintf("%s:%d", host, port)
    fmt.Println("serving at " + addr)
    http.ListenAndServe(addr, nil)
}
