package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var Logger *AsyncLogger

type AsyncLogger struct {
	log        *log.Logger
	logChannel chan string
	wg         sync.WaitGroup
	closed     chan struct{}
	closing    chan struct{}
}

func Initialise() {
	out := os.Stdout
	bufferSize := 100
	prefix := ""
	flags := log.LstdFlags

	logger := &AsyncLogger{
		log:        log.New(out, prefix, flags),
		logChannel: make(chan string, bufferSize),
		closed:     make(chan struct{}),
		closing:    make(chan struct{}),
	}

	logger.wg.Add(1)
	go logger.flusher()

	Logger = logger
}

// flushes log messages in the background
func (logger *AsyncLogger) flusher() {
	defer logger.wg.Done()
	for {
		select {
		case msg := <-logger.logChannel:
			logger.log.Println(msg)
		case <-logger.closing:
			for {
				select {
				case msg := <-logger.logChannel:
					logger.log.Println(msg)
				default:
					close(logger.closed)
					return
				}
			}
		}
	}
}

func (logger *AsyncLogger) Println(v ...any) {
	// remove trailing newline
	message := fmt.Sprintln(v...)
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}
	select {
	case logger.logChannel <- message:
	default:
		// drop if buffer is full
	}
}

func (logger *AsyncLogger) Printf(format string, v ...any) {
	select {
	case logger.logChannel <- fmt.Sprintf(format, v...):
	default:
		// drop if buffer is full
	}
}

func (logger *AsyncLogger) Shutdown() {
	close(logger.closing)
	logger.wg.Wait()
}

// wrapper functions for export
func Println(v ...any) {
	if Logger != nil {
		Logger.Println(v...)
	}
}

func Printf(format string, v ...any) {
	if Logger != nil {
		Logger.Printf(format, v...)
	}
}

func Shutdown() {
	if Logger != nil {
		Logger.Shutdown()
	}
}
