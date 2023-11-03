package main

import (
	"flag"
	"fmt"
	"time"

	// "github.com/yutfut/log-streamer/ch"
	"github.com/yutfut/log-streamer/watcher"
	"github.com/yutfut/log-streamer/writer"
	"github.com/yutfut/log-streamer/yetSender"
	"github.com/yutfut/log-streamer/tools"
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
	allFiles := make([]string, 0)

	flag.Var(&files, "file", "go run ./main.go -file file.log -file file1.log -file file2.log")
	flag.Parse()

	fmt.Println(files)

	allFiles = append(allFiles, files...)

	config := tools.NewConfig()

	err := tools.ReadConfigFile("conf.toml", config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)

	allFiles = append(allFiles, config.Files...)

	fmt.Println(allFiles)

	// watcher := watcher.NewWatcher(ch.NewClickHouse())
	watcher := watcher.NewWatcher(yetSender.NewSender())

	writer := writer.NewWatcher()

	writer.AddFiles(allFiles)
	go writer.Start()

	watcher.AddFiles(allFiles)
	go watcher.Start()

	time.Sleep(10 * time.Second)
	fmt.Println("STOP")

	writer.Stop()
	watcher.Stop()

	time.Sleep(2 * time.Second)
}
