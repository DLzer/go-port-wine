package main

import (
	"time"

	"github.com/DLzer/go-port-wine/port"
	"golang.org/x/sync/semaphore"
)

func main() {

	ps := &port.PortScanner{
		hostname: "127.0.0.1",
		lock:     semaphore.NewWeighted(port.Ulimit()),
	}

	ps.Start(1, 65535, 500*time.Millisecond)

}
