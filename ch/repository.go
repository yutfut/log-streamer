package ch

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse struct {
    driver  driver.Conn
}

func NewClickHouse(driver driver.Conn) *ClickHouse {
    return &ClickHouse{
        driver: driver,
    }
}

func (CH *ClickHouse) InsertLog(log string) {
    ctx := context.Background()
    _, err := CH.driver.Query(ctx, "INSERT INTO logs(log, timestamp) VALUES ($1, now())", log)
    if err != nil {
        fmt.Println(err)
    }
}