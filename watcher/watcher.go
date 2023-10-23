package watcher

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/yutfut/log-streamer/ch"
)

type WatcherInterface interface {
	AddFiles(files []string) error
	Start() error
	Stop() error
}

type Watcher struct {
	files      []string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewWatcher() *Watcher {
	ctxBuff, cancelFuncBuff := context.WithCancel(context.Background())
	return &Watcher{
		ctx:        ctxBuff,
		cancelFunc: cancelFuncBuff,
	}
}

func (w *Watcher) AddFiles(files []string) error {
	if len(files) == 0 {
		return errors.Errorf("files is empty")
	}

	w.files = append(w.files, files...)

	return nil
}

func readerLog(wg *sync.WaitGroup, file *os.File, che *ch.ClickHouse, filePath string, ctx context.Context) {
// func readerLog(wg *sync.WaitGroup, file *os.File, filePath string, ctx context.Context) {
	defer wg.Done()

	reader := bufio.NewReader(file)

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    watcher.Add(filePath)

	for {
		line, _, err := reader.ReadLine()
        if err != nil {
            if err == io.EOF {
                fmt.Println(err, filePath)

                select {
                case event, ok := <-watcher.Events:
                    if !ok {
                        return
                    }
                    fmt.Println(event)
                    continue
                case err, ok := <-watcher.Errors:
                    if !ok {
                        return
                    }
                    fmt.Println(err, filePath)
                }
            } else {
                fmt.Println(err)
            }
        }

		if len(line) != 0 {
			fmt.Println(string(line))
			che.InsertLog(string(line), filePath)
		}
        
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (w *Watcher) Start() error {
	conn, err := ch.Connect()
	if err != nil {
		panic((err))
	}

	clickHouseDriver := ch.NewClickHouse(conn)

	// tike out in WatcherInterface struct
	wg := &sync.WaitGroup{}

	for _, file := range w.files {

		fileOutput, err := os.Open(file)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
			return err
		}
		defer fileOutput.Close()

		wg.Add(1)

		go readerLog(wg, fileOutput, clickHouseDriver, file, w.ctx)
        // go readerLog(wg, fileOutput, file, w.ctx)
	}

	wg.Wait()
	return nil
}

func (w *Watcher) Stop() {
	w.cancelFunc()
}
