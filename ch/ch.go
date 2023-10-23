package ch

type CH interface {
	Sender(log, file string) error
}
