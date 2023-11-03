package writer

import (
	"context"
	"log"
	"os"
	"time"

	"sync"

	"github.com/pkg/errors"
)

func writeLog(wg *sync.WaitGroup, file *os.File, ctx context.Context) {
	defer wg.Done()

	logger := log.New(file, "", log.Ltime)
	logger.Println("start")

	var i uint64 = 0
	for {
		i += 1
		logger.Println(i)
		time.Sleep(5 * time.Second)

		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

type WriterInterface interface {
	AddFiles(files []string) error
	Start() error
	Stop() error
}

type Writer struct {
	files      []string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewWatcher() *Writer {
	ctxBuff, cancelFuncBuff := context.WithCancel(context.Background())
	return &Writer{
		ctx:        ctxBuff,
		cancelFunc: cancelFuncBuff,
	}
}

func (w *Writer) AddFiles(files []string) error {
	if len(files) == 0 {
		return errors.Errorf("files is empty")
	}

	w.files = append(w.files, files...)

	return nil
}

func (w *Writer) Start() error {
	// tike out in WatcherInterface struct
	wg := &sync.WaitGroup{}

	for _, file := range w.files {

		// fileInput, err := os.Create(file)
		// if err != nil {
		// 	log.Println("error create file:", err)
		// 	return err
		// }

		fileInput, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("writer error opening file:", err)
			return err
		}
		defer fileInput.Close()

		wg.Add(1)

		go writeLog(wg, fileInput, w.ctx)
	}

	wg.Wait()
	return nil
}

func (w *Writer) Stop() {
	w.cancelFunc()
}
