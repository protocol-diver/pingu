package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/protocol-diver/pingu"
)

// go run ping.go --timeout 1000 127.0.0.1:1111
func main() {
	timeoutFlag := flag.Int("timeout", 1000, "milliseconds")
	flag.Parse()

	dest := (flag.Arg(0))
	timeout := time.Duration(*timeoutFlag * 1000000)

	p, err := pingu.NewPingu("127.0.0.1:3821", nil)
	if err != nil {
		return
	}

	p.Start()
	defer p.Close()

	err = p.PingPongWithRawAddr(dest, timeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ping-pong success ip:", dest, "timeout", timeout)
}
