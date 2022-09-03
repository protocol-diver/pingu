package main

import (
	"fmt"
	"net"
	"time"

	"github.com/dbadoy/pingu"
)

func main() {
	client, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8753})
	p := pingu.NewPingu(client, pingu.Config{})
	p.Start()

	// peer2
	p.Register("127.0.0.1:8754")
	ticker := time.NewTicker(5 * time.Second)
	_ = p.BroadcastPingWithTicker(*ticker, 5*time.Second)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(p.PingTable())
	}
}
