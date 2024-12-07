package main

import (
	"flag"
	"log"
	"time"

	"github.com/realrabbithouse/go-play/filelock"
)

var (
	hold    time.Duration
	timeout time.Duration
	shared  bool
	block   bool
)

func init() {
	flag.DurationVar(&hold, "h", 10*time.Second, "hold duration in seconds")
	flag.DurationVar(&timeout, "t", 5*time.Second, "timeout waiting for lock")
	flag.BoolVar(&shared, "s", false, "use as shared lock or not")
	flag.BoolVar(&block, "b", false, "wait for lock to be acquired")
}

func main() {
	flag.Parse()

	l, err := filelock.New("example.lock")
	if err != nil {
		log.Println("filelock.New failed:", err)
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			log.Println("Unlock failed:", err)
		}
	}()

	var opts []filelock.Option
	if block {
		opts = append(opts, filelock.WithBlock())
	}
	opts = append(opts, filelock.WithTimeout(timeout), filelock.WithRemove())

	if shared {
		if err := l.RLock(opts...); err != nil {
			log.Println("RLock failed:", err)
			return
		}
	} else {
		if err := l.WLock(opts...); err != nil {
			log.Println("WLock failed:", err)
			return
		}
	}

	// Perform critical section.
	log.Println("Locked for", hold, "seconds")
	time.Sleep(hold)
	log.Println("Done")
}
