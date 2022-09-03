package main

import (
	"fmt"
	"net"
	"pingu"
	"time"
)

func main() {
	client, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8753})

	p := pingu.NewPingu(client, pingu.Config{RecvBufferSize: 256})
	p.Start()

	p.Register("127.0.0.1:8754")
	p.Register("127.0.0.1:8755")
	p.Register("127.0.0.1:8756")

	ticker := time.NewTicker(3 * time.Second)
	_ = p.BroadcastPingWithTicker(*ticker, 3*time.Second)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(p.PingTable())

		time.Sleep(1 * time.Second)
		p.Stop()

		time.Sleep(1 * time.Second)
		p.Start()

	}
}
