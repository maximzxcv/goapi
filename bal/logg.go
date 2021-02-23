package bal

import "fmt"

// Logger ....
type Logger interface {
	Debug(a ...interface{})
	Error(error error)
}

// Logg ...
type Logg struct{}

// NewLogg ...
func NewLogg() *Logg {
	return &Logg{}
}

// Debug ...
func (l *Logg) Debug(a ...interface{}) {
	fmt.Println(a)
}

// Error ...
func (l *Logg) Error(err error) {
	fmt.Println(err)
}
