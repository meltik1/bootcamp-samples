package Cache

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rd *redis.Client
}

func (cache Cache) Init(ctx context.Context, rows int) error {
	for i := 0; i < rows; i++ {
		err := cache.rd.Set(ctx, fmt.Sprintf("user_%d", i), "robert paulson", 0).Err()
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
	}
	return nil
}

func New(client *redis.Client) Cache {
	return Cache{rd: client}
}

func (cache Cache) Request(ctx context.Context, count int) {
	for i := 0; i < count; i++ {
		id := rand.Int() % 1000
		res, err := cache.rd.Get(ctx, fmt.Sprintf("user_%d", id)).Result()
		if err != nil {
			fmt.Println("Error occured while making request")
			return
		}
		fmt.Printf("%s", res)
	}
}

func (cache Cache) Hash(ctx context.Context, smth string) error {
	set := cache.rd.HSet(ctx, "LOL", smth)
	if set.Err() != nil {
		return errors.New("Error in hashtale")
	}
	return nil
}
