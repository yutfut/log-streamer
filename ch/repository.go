package ch

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type CH interface {
	Sender(log, file string) error
}

type ClickHouse struct {
	driver driver.Conn
}

func NewClickHouse() *ClickHouse {
    conn, err := Connect()
	if err != nil {
		panic((err))
	}

	return &ClickHouse{
		driver: conn,
	}
}

func (ch *ClickHouse) Sender(log, file string) error {
	ctx := context.Background()

	var args []interface{}
	args = append(args, log)
	args = append(args, file)

	_, err := ch.driver.Query(ctx, "INSERT INTO logs(id, log, file, timestamp) VALUES (generateUUIDv4(), $1, $2, now())", args...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
