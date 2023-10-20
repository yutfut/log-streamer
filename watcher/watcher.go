package watcher

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"log"

	"github.com/yutfut/log-streamer/ch"
    "sync"

	"github.com/pkg/errors"
)

type Watcher interface {
	AddFiles(files []string) error
    Start() error
    // Stop() error
}

type WatcherS struct {
	files []string
}

func (w *WatcherS) AddFiles(files []string) error {
	if len(files) == 0 {
		return errors.Errorf("files is empty")
	}

	for _, file := range files {
		w.files = append(w.files, file)
	}

	return nil
}

func readerLog(wg *sync.WaitGroup, file *os.File, che *ch.ClickHouse, fileName string) {
    defer wg.Done()
    
    reader := bufio.NewReader(file)

    for {
        line, _, err := reader.ReadLine()
        if err != nil && err != io.EOF {
            fmt.Println(err)
        }

        if len(line) != 0 {
            fmt.Println(string(line))
            che.InsertLog(string(line), fileName)
        }
    }
}

func (w *WatcherS) Start() error {
	conn, err := ch.Connect()
    if err != nil {
        panic((err))
    }

    clickHouseDriver := ch.NewClickHouse(conn)

	// tike out in Watcher struct
	wg := &sync.WaitGroup{}

    for _, file := range w.files {
    
        fileOutput, err := os.Open(file)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
			return err
        }
        defer fileOutput.Close()

        wg.Add(1)
    
        go readerLog(wg, fileOutput, clickHouseDriver, file)
    }

	wg.Wait()
	return nil
}