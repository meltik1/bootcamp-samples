package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/devhands-io/bootcamp-samples/golang/vanilla/Cache"
	"github.com/devhands-io/bootcamp-samples/golang/vanilla/handlers"
	"github.com/devhands-io/bootcamp-samples/golang/vanilla/payload"
	"github.com/devhands-io/bootcamp-samples/golang/vanilla/storage"
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

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	database := storage.New(conn)

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong, err)
		return
	}

	cache := Cache.New(rdb)

	runtime.GOMAXPROCS(2 * runtime.NumCPU())

	flag.Parse()

	// dummy handlers
	http.HandleFunc("/", handlers.Ok)
	http.HandleFunc("/hello", handlers.Hello)

	// payload
	cpuSleep := payload.NewGetrusagePayload()
	ioSleep := payload.NewIOPayload()
	http.HandleFunc("/payload", handlers.SleepHandler(cpuSleep, ioSleep, database, cache))
	http.HandleFunc("/init", handlers.InitDBHandler(database, cache))

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("serving at " + addr)
	http.ListenAndServe(addr, nil)
}
