package main

import (
    "fmt"
    "log"
    "time"
    "os"
    "bufio"

    "github.com/yutfut/log-streamer/ch"
    "sync"
)

func writeLog(wg *sync.WaitGroup, file *os.File) {
    defer wg.Done()

	logger := log.New(file, "", log.Ltime)
    logger.Println("start")

    var i uint64 = 0
    for {
        i+=1
        logger.Println(i)
        time.Sleep(5 * time.Second)
    }
}

func readerLog(wg *sync.WaitGroup, file *os.File, che *ch.ClickHouse) {
    defer wg.Done()
    
    reader := bufio.NewReader(file)

    for {
        line, prefix, err := reader.ReadLine()
        fmt.Println(prefix)
        if err != nil {
            fmt.Println(err)
        }
        if len(line) != 0 {
            fmt.Println(string(line))
            che.InsertLog(string(line))
        } else {
            fmt.Println("--//--")
            time.Sleep(time.Second)
        }
    }
}

func main() {
    conn, err := ch.Connect()
    if err != nil {
        panic((err))
    }

    che := ch.NewClickHouse(conn)

    f, err := os.OpenFile("file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

    wg := &sync.WaitGroup{}
	fmt.Println("hello")

    wg.Add(2)
    go writeLog(wg, f)

    go readerLog(wg, f, che)
    wg.Wait()
}