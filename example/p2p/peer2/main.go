package main

import (
	"fmt"
	"net"
	"pingu"
	"time"
)

func main() {
	client, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8754})
	p := pingu.NewPingu(client, pingu.Config{})
	p.Start()

	// peer1
	p.Register("127.0.0.1:8753")
	ticker := time.NewTicker(5 * time.Second)
	_ = p.BroadcastPingWithTicker(*ticker, 5*time.Second)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(p.PingTable())
	}
}
