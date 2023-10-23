package main

import (
	"flag"
	"fmt"
	"time"

	// "github.com/yutfut/log-streamer/ch"
	"github.com/yutfut/log-streamer/watcher"
	"github.com/yutfut/log-streamer/writer"
	"github.com/yutfut/log-streamer/yetsender"
)

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

	// watcher := watcher.NewWatcher(ch.NewClickHouse())
	watcher := watcher.NewWatcher(yetsender.NewSender())

	writer := writer.NewWatcher()

	writer.AddFiles(files)
	go writer.Start()

	watcher.AddFiles(files)
	go watcher.Start()

	time.Sleep(10 * time.Second)
	fmt.Println("STOP")

	writer.Stop()
	watcher.Stop()

	time.Sleep(2 * time.Second)
}
