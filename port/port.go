package port

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	hostname string
	lock     *semaphore.Weighted
}

type ScanResult struct {
	Port  int
	State string
}

func ScanPort(hostname string, port int, timeout time.Duration) ScanResult {

	result := ScanResult{Port: port}

	target := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
		} else {
			result.State = "Closed"
			return result
		}
	}

	conn.Close()
	result.State = "Open"
	return result
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	waitGroup := sync.WaitGroup{}
	defer waitGroup.Wait()

	for port := f; port <= l; port++ {
		waitGroup.Add(1)
		ps.lock.Acquire(context.TODO(), 1)

		go func(port int) {
			defer ps.lock.Release(1)
			defer waitGroup.Done()
			ScanPort(ps.hostname, port, timeout)
		}(port)
	}
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}

	return i
}
