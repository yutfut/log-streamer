package main

import (
    "fmt"
    "flag"
    "log"
    "time"
    "os"

    "github.com/yutfut/log-streamer/watcher"
    "sync"
)

// type Watcher interface {
//     Start() error
//     Stop() error
// }
    
// type Sender interface {
//     Send(message string)
// }

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

// func readerLog(wg *sync.WaitGroup, file *os.File, che *ch.ClickHouse, fileName string) {
//     defer wg.Done()
    
//     reader := bufio.NewReader(file)

//     for {
//         line, _, err := reader.ReadLine()
//         if err != nil && err != io.EOF {
//             fmt.Println(err)
//         }

//         if len(line) != 0 {
//             fmt.Println(string(line))
//             che.InsertLog(string(line), fileName)
//         }
//     }
// }


// https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang

type arrayFlags []string

func (i *arrayFlags) String() string {
    return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}

var files arrayFlags

func main() {
    flag.Var(&files, "file", "go run ./main.go -file file.log -file file1.log -file file2.log")
    flag.Parse()

    fmt.Println(files)

    w := watcher.NewWatcher()

    w.AddFiles(files)
    w.Start()

    wg := &sync.WaitGroup{}

    for _, file := range files {

        fileInput, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
        }
        defer fileInput.Close()

        wg.Add(1)
    
        go writeLog(wg, fileInput)
    }

    wg.Wait()

    // conn, err := ch.Connect()
    // if err != nil {
    //     panic((err))
    // }

    // clickHouseDriver := ch.NewClickHouse(conn)

    // wg := &sync.WaitGroup{}

    // for _, file := range files {

    //     fileInput, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    //     if err != nil {
    //         log.Fatalf("error opening file: %v", err)
    //     }
    //     defer fileInput.Close()
    
    //     fileOutput, err := os.Open(file)
    //     if err != nil {
    //         log.Fatalf("error opening file: %v", err)
    //     }
    //     defer fileOutput.Close()

    //     wg.Add(2)
    
    //     go writeLog(wg, fileInput)
    
    //     go readerLog(wg, fileOutput, clickHouseDriver, file)
    // }

    // wg.Wait()
}