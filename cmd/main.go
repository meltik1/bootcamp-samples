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

	"github.com/meltik/study/internal/infra/cache"
	"github.com/meltik/study/internal/infra/handlers/init"
	"github.com/meltik/study/internal/infra/storage"
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

	err := godotenv.Load("/Users/nszuev/GolandProjects/bootcamp-samples/.env.paas")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	database := storage.New(conn)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_LOGIN"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong, err)
		return
	}

	cache := cache.New(rdb)

	runtime.GOMAXPROCS(2 * runtime.NumCPU())

	flag.Parse()
	http.HandleFunc("/init", init.InitDBHandler(database, cache))

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("serving at " + addr)
	http.ListenAndServe(addr, nil)
}
