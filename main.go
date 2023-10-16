package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "net"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)



func main() {
	fmt.Println("hello")

	conn, err := connection1()
    if err != nil {
        panic((err))
    }

    fmt.Println(1)

	ctx := context.Background()
    rows, err := conn.Query(ctx, "SELECT * FROM my_user")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(2)

	for rows.Next() {
        var (
            user_id uint32
            name string
        )
        if err := rows.Scan(
            &user_id,
            &name,
        ); err != nil {
            log.Fatal(err)
        }
        fmt.Println("user_id: ", user_id, " name: ", name)
    }

    fmt.Println(3)
}

func connection1() (driver.Conn, error) {
    dialCount := 0
    conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      time.Second * 30,
		MaxOpenConns:     5,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize: 10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "log-streamer", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return nil, err
	}

    return conn, nil
}