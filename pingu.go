// Copyright (c) 2022, Seungbae Yu <dbadoy4874@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pingu

import (
	"fmt"
	"net"
	"net/netip"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PingType = 1 + iota
	// NotificationType

	MaxPacketSize = 4

	localhost   = "127.0.0.1"
	defaultPort = 4874
)

type Pingu struct {
	conn *net.UDPConn
	cfg  *Config

	// Pingu has full open about ping request. If recv ping, always send pong.
	// But if recv pong from not exist in white list, it won't store even if
	// send a ping request self.
	wl map[string]bool

	// 'peers' mapping rawAddress to health status.
	// The health status set when the ping-pong request completes
	peers map[string]bool

	recvPongs chan Packet

	isRun uint32
	mu    sync.Mutex

	stop chan struct{}
}

func DefaultAddress() string {
	return fmt.Sprintf("%s:%d", localhost, defaultPort)
}

// Pingu is not accept double use to net.UDPConn. It's should only be used once.
// For avoid confuse, generate net.UDPConn in NewPingu.
func NewPingu(rawAddr string, cfg *Config) (*Pingu, error) {
	conn, err := listenWithRawAddr(rawAddr)
	if err != nil {
		return nil, err
	}
	// Works if succed generate net.UDPConn.
	if cfg == nil {
		cfg = new(Config)
		cfg.Default()
	}
	if cfg.RecvBufferSize < 1 {
		cfg.RecvBufferSize = 256
	}
	return &Pingu{
		conn:      conn,
		cfg:       cfg,
		wl:        make(map[string]bool),
		peers:     make(map[string]bool),
		stop:      make(chan struct{}, 1),
		recvPongs: make(chan Packet, cfg.RecvBufferSize),
	}, nil
}

func (p *Pingu) Start() {
	if atomic.LoadUint32(&p.isRun) == 1 {
		return
	}
	atomic.StoreUint32(&p.isRun, 1)
	go p.detectLoop()
}

func (p *Pingu) Stop() {
	if atomic.LoadUint32(&p.isRun) == 0 {
		return
	}
	p.peers = make(map[string]bool)
	atomic.StoreUint32(&p.isRun, 0)
	p.stop <- struct{}{}
}

func (p *Pingu) detectLoop() {
	for {
		select {
		case <-p.stop:
			if p.cfg.Verbose {
				fmt.Println("[pingu] recv close signal")
			}
			return
		default:
			b := make([]byte, MaxPacketSize)
			size, sender, err := p.conn.ReadFromUDP(b)
			if size == 0 {
				continue
			}
			if err != nil {
				if p.cfg.Verbose {
					fmt.Printf("[pingu] ReadFromUDP error %v", err)
				}
				continue
			}

			// Set sender when before start goroutine.
			// Not after started goroutine. It may not thread safety.
			go func(sender *net.UDPAddr) {
				packet, err := ParsePacket(b, sender)
				if err != nil {
					if p.cfg.Verbose {
						fmt.Printf("[pingu] detected invalid protocol, reason : %v\n", err)
					}
					return
				}
				switch packet.Kind() {
				case Ping:
					go p.pong([]*net.UDPAddr{sender})
				case Pong:
					p.recvPongs <- packet
				default:
					panic(fmt.Sprintf("[pingu] detected invalid protocol: invalid packet type %v", packet.Kind()))
				}
			}(sender)
		}
	}
}

func (p *Pingu) RemoteAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *Pingu) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

// Register is register to broadcast list that input address.
func (p *Pingu) Register(addr *net.UDPAddr) {
	p.register(addr.String())
}

func (p *Pingu) RegisterWithRawAddr(raw string) error {
	_, err := rawAddrToUDPAddr(raw)
	if err != nil {
		return err
	}
	p.register(raw)
	return nil
}

// Unregister is remove input address from broadcast list.
func (p *Pingu) Unregister(addr *net.UDPAddr) {
	p.unregister(addr.String())
}

func (p *Pingu) UnregisterWithRawAddr(raw string) error {
	_, err := rawAddrToUDPAddr(raw)
	if err != nil {
		return err
	}
	p.unregister(raw)
	return nil
}

// PingPong sends a 'ping' and waits for a 'pong' to be received.
func (p *Pingu) PingPong(addr *net.UDPAddr, timeout time.Duration) error {
	return p.pingpong(addr, timeout)
}

func (p *Pingu) PingPongWithRawAddr(raw string, timeout time.Duration) error {
	addr, err := rawAddrToUDPAddr(raw)
	if err != nil {
		return err
	}
	return p.pingpong(addr, timeout)
}

func (p *Pingu) register(rawAddr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.wl[rawAddr] = true
}

func (p *Pingu) unregister(rawAddr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.wl[rawAddr] = false
}

func (p *Pingu) pingpong(addr *net.UDPAddr, timeout time.Duration) error {
	rawAddr := addr.String()
	p.mu.Lock()
	registered := p.wl[rawAddr]
	p.mu.Unlock()

	if !registered {
		return fmt.Errorf("not registered ip: %v" + rawAddr)
	}

	p.ping([]*net.UDPAddr{addr}, timeout)
	if !p.IsAlive(rawAddr) {
		return fmt.Errorf("ping-pong failed ip: %v, timeout: %v", rawAddr, timeout)
	}
	return nil
}

// Send broadcast with ticker.
func (p *Pingu) BroadcastPingWithTicker(ticker time.Ticker, per time.Duration) chan struct{} {
	var cancel chan struct{}
	go func() {
		for {
			select {
			case <-ticker.C:
				// If 'per' greater than ticker duration, ticker wait broadcasePing done.
				// Do not call broadcastPing by goroutine. If you use goroutine, will accumulate
				// meaningless running goroutines.
				p.broadcastPing(per)
			case <-cancel:
				return
			}
		}
	}()
	return cancel
}

// Do not call by goroutine. It's running it once is enough.
func (p *Pingu) broadcastPing(timeout time.Duration) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	go p.broadcast(PingType, timeout)
}

func (p *Pingu) IsAlive(raw string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peers[raw]
}

func (p *Pingu) PingTable() map[string]bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.snapPingTable()
}

// snapPingTable returns deep copied map about broadcast list.
//
// The caller must hold b.mu.
func (p *Pingu) snapPingTable() (r map[string]bool) {
	r = make(map[string]bool, len(p.peers))
	for addr, health := range p.peers {
		r[addr] = health
	}
	return
}

func (p *Pingu) broadcast(t byte, timeout time.Duration) {
	p.mu.Lock()
	addrs := make([]*net.UDPAddr, 0, len(p.wl))
	for target := range p.wl {
		addrs = append(addrs, mustAddrToUDPAddr(target))
	}
	p.mu.Unlock()
	if len(addrs) == 0 {
		if p.cfg.Verbose {
			fmt.Println("[pingu] there is no target")
		}
		return
	}
	switch t {
	case PingType:
		p.ping(addrs, timeout)
	default:
		panic(fmt.Sprintf("[pingu] detected invalid protocol: invalid packet type %v", t))
	}
}

func (p *Pingu) ping(addrs []*net.UDPAddr, timeout time.Duration) {
	for _, addr := range addrs {
		packet := new(PingPacket)
		packet.SetKind(Ping)
		byt, _ := SuitablePack(packet)

		if _, err := p.conn.WriteToUDP(byt, addr); err != nil {
			fmt.Println(err)
			break
		}
	}

	// The snapshot that marking the changed peer status. If get 'pong',
	// remove sender from snapshot. This means that peers that did not
	// send a response to the PING remain in the snapshot.
	p.mu.Lock()
	tempSnapTable := p.snapPingTable()
	p.mu.Unlock()

	// This is the case that not requesting a heartbeat for all peers.
	// For update only requested peers.
	if len(addrs) != len(tempSnapTable) {
		t := make(map[string]bool, len(addrs))
		for _, addr := range addrs {
			rawAddr := addr.String()
			t[rawAddr] = tempSnapTable[rawAddr]
		}
		tempSnapTable = t
	}

	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			p.mu.Lock()
			defer p.mu.Unlock()
			for rawAdrr := range tempSnapTable {
				p.peers[rawAdrr] = false
			}
			if p.cfg.Verbose {
				fmt.Println("[pingu] ", p.snapPingTable())
			}
			return
		case r := <-p.recvPongs:
			rawAddr := (*r.Sender()).String()
			p.mu.Lock()
			if p.wl[rawAddr] {
				p.peers[rawAddr] = true
			}
			p.mu.Unlock()
			delete(tempSnapTable, rawAddr)
		}
	}
}

func (p *Pingu) pong(addrs []*net.UDPAddr) {
	for _, addr := range addrs {
		packet := new(PongPacket)
		packet.SetKind(Pong)
		byt, _ := SuitablePack(packet)

		if _, err := p.conn.WriteToUDP(byt, addr); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

// [Benchmark]
//		net.ResolveUDPAddr					10000000                151.0 ns/op
//		netip.MustParseAddrPort, net.UDPAddrFromAddrPort	10000000                62.55 ns/op
func rawAddrToUDPAddr(s string) (*net.UDPAddr, error) {
	rawAddr, err := netip.ParseAddrPort(s)
	if err != nil {
		return nil, err
	}
	return net.UDPAddrFromAddrPort(rawAddr), nil
}

func mustAddrToUDPAddr(s string) *net.UDPAddr {
	rawAddr := netip.MustParseAddrPort(s)
	return net.UDPAddrFromAddrPort(rawAddr)
}

func listenWithRawAddr(rawAddr string) (*net.UDPConn, error) {
	addr, err := rawAddrToUDPAddr(rawAddr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}
