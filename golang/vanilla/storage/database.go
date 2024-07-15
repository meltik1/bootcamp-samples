package storage

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5"

	"github.com/devhands-io/bootcamp-samples/golang/vanilla/models"
)

type DB struct {
	connect *pgx.Conn
}

func New(connect *pgx.Conn) DB {
	return DB{connect: connect}
}

func (d DB) Request(ctx context.Context, count int) {
	var u models.User
	for i := 0; i < count; i++ {
		id := rand.Int() % 1000
		err := d.connect.QueryRow(ctx, "select * from Users u where u.id = $1", id).Scan(u)
		if err != nil {
			fmt.Println("Error occured while making request")
			return
		}
		fmt.Printf("%v", u)
	}
}

func (d DB) Init(ctx context.Context, rows int) error {
	for i := 0; i < rows; i++ {
		_, err := d.connect.Exec(ctx, "insert into Users(id, name, surname) values ($1, $2, 'paulson') on conflict do nothing ", i, fmt.Sprintf("robert_%d", i))
		if err != nil {
			return err
		}
	}
	return nil
}
