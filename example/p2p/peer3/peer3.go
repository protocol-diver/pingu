package main

import (
	"fmt"
	"time"

	"github.com/dbadoy/pingu"
)

func main() {
	p, err := pingu.NewPingu("127.0.0.1:8755", nil)
	if err != nil {
		return
	}
	p.Start()

	// peer1
	p.Register("127.0.0.1:8753")
	// peer2
	p.Register("127.0.0.1:8754")
	ticker := time.NewTicker(5 * time.Second)
	_ = p.BroadcastPingWithTicker(*ticker, 5*time.Second)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(p.PingTable())
	}
}
