// DO NOT USE THIS

package main

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"sync/atomic"
	"time"

	"github.com/protocol-diver/pingu"
)

func main() {
	var perSleep int
	flag.IntVar(&perSleep, "rest_time", 3, "")
	fmt.Printf("rest time set: %v\n", perSleep)

	tpingu, err := pingu.NewPingu("127.0.0.1:8771", &pingu.Config{RecvBufferSize: 12})
	if err != nil {
		fmt.Println(err)
		return
	}
	client, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:8770")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		client.Close()
		// tpingu.Close()
	}()

	tpingu.Start()

	count := 5000
	increase := 1000
	for j := 0; j < 10; j++ {
		time.Sleep(time.Duration(perSleep) * time.Second)
		in := make(chan struct{})
		out := make(chan struct{})
		for i := 0; i < count; i++ {
			go func() {
				ping := new(pingu.PingPacket)
				b, _ := pingu.SuitableUnpack(ping)
				if _, err := client.WriteToUDP(b, net.UDPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:8771"))); err != nil {
					return
				}
				in <- struct{}{}

				response := make([]byte, 4)
				client.ReadFrom(response)
				var pongPacket pingu.PongPacket
				if err := pingu.SuitablePack(response, &pongPacket); err != nil {
					return
				}
				out <- struct{}{}
			}()
			count += increase
		}

		var incount int32
		var outcount int32
		start := time.Now()
	T:
		for {
			select {
			case <-in:
				atomic.AddInt32(&incount, 1)
				if atomic.LoadInt32(&incount) == atomic.LoadInt32(&outcount) {
					break T
				}
			case <-out:
				atomic.AddInt32(&outcount, 1)
				if atomic.LoadInt32(&incount) == atomic.LoadInt32(&outcount) {
					break T
				}
			}
		}
		end := time.Since(start)
		fmt.Printf("sum: %d, total in: %d, out: %d, spent time: %v\n", count, incount, outcount, end)
	}
}
