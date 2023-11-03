package watcher

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type WatcherInterface interface {
	AddFiles(files []string) error
	Start() error
	Stop() error
}

type senderInterface interface {
	Sender(log, file string) error
}

type Watcher struct {
	sender     senderInterface
	files      []string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewWatcher(sender senderInterface) *Watcher {
	ctxBuff, cancelFuncBuff := context.WithCancel(context.Background())
	return &Watcher{
		sender:     sender,
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

func readerLog(wg *sync.WaitGroup, file *os.File, che senderInterface, filePath string, ctx context.Context) {
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
			che.Sender(string(line), filePath)
		}

		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func waiterFile(file string) (*os.File, error) {
	var err error
	for i := 0; i < 10; i++ {
		fileOutput, err := os.Open(file)
		if err == nil {
			return fileOutput, nil
		}
		time.Sleep(time.Second)
	}
	log.Fatalf("watcher error opening file: %v", err)
	return nil, err
}

func (w *Watcher) Start() error {
	// tike out in WatcherInterface struct
	wg := &sync.WaitGroup{}

	for _, file := range w.files {

		fileOutput, err := waiterFile(file)
		if err != nil {
			log.Fatalf("watcher waiter file: %v", err)
			return err
		}

		defer fileOutput.Close()

		wg.Add(1)

		go readerLog(wg, fileOutput, w.sender, file, w.ctx)
		// go readerLog(wg, fileOutput, file, w.ctx)
	}

	wg.Wait()
	return nil
}

func (w *Watcher) Stop() {
	w.cancelFunc()
}
