package logging

import (
	"fmt"
	"log"

	"go.uber.org/fx"
)

// Logger is a simple logging interface
type Logging struct{}

func Logger() *Logging {
	return &Logging{}
}

// Logging events.
func (l *Logging) Event(d ...any) {
	name := "logging:"
	s := fmt.Sprint(d...)
	message := fmt.Sprintf("%s %s", name, s)

	log.Println(message)
}

// Module groups the Logger constructor for Fx
var Module = fx.Provide(Logger)
