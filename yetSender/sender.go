package yetSender

import "fmt"

type Sender struct {}

func NewSender() *Sender{
	return &Sender{}
}

func (Sender) Sender(log, file string) error {
	fmt.Println(log, file)
	return nil
}