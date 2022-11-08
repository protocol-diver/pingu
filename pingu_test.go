package pingu_test

import (
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/protocol-diver/pingu"
)

func TestNewPinguAddr(t *testing.T) {
	type td struct {
		got    string
		expect pingu.Config
		err    string
	}
	tdl := []td{
		{got: "127.0.0.1:1723", expect: pingu.Config{}, err: ""},
		{got: "127.0.0.1:", expect: pingu.Config{}, err: "no port"},
		{got: ":3821", expect: pingu.Config{}, err: "no IP"},
		{got: "", expect: pingu.Config{}, err: "not an ip:port"},
		{got: "hello", expect: pingu.Config{}, err: "not an ip:port"},
	}

	for _, td := range tdl {
		_, err := pingu.NewPingu(td.got, nil)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("NewPingu failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
	}
}

func TestPingPong(t *testing.T) {
	pingu1, err := pingu.NewPingu("127.0.0.1:9190", nil)
	if err != nil {
		t.Fatalf("PingPong NewPingu failure %v", err)
	}
	defer pingu1.Close()
	pingu2, err := pingu.NewPingu("127.0.0.1:9191", nil)
	if err != nil {
		t.Fatalf("PingPong NewPingu failure %v", err)
	}
	defer pingu2.Close()

	pingu1.Start()
	pingu2.Start()

	// PingPongWithRawAddr
	err = pingu1.PingPongWithRawAddr("127.0.0.1:9191", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("PingPong failure got: %v", err)
	}
	// pingpong to unknown pingu
	err = pingu1.PingPongWithRawAddr("127.0.0.1:9195", 100*time.Millisecond)
	if err == nil {
		t.Fatalf("PingPong failure got: %v", err)
	}
	// check the duration
	now := time.Now()
	err = pingu1.PingPongWithRawAddr("127.0.0.1:9195", 100*time.Millisecond)
	if err == nil {
		t.Fatalf("PingPong failure got: %v", err)
	}
	spent := time.Since(now)
	if spent < 100*time.Millisecond {
		t.Fatalf("PingPong failure: incorrect act about pingpong timeout")
	}

	// PingPong
	rawAddr := netip.MustParseAddrPort("127.0.0.1:9191")
	err = pingu1.PingPong(net.UDPAddrFromAddrPort(rawAddr), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("PingPong failure got: %v", err)
	}
	// pingpong to unknown pingu
	rawAddr = netip.MustParseAddrPort("127.0.0.1:9195")
	err = pingu1.PingPong(net.UDPAddrFromAddrPort(rawAddr), 100*time.Millisecond)
	if err == nil {
		t.Fatalf("PingPong failure got: %v", err)
	}
}

func TestRegister(t *testing.T) {
	pingu1, err := pingu.NewPingu("127.0.0.1:9190", nil)
	if err != nil {
		t.Fatalf("Register NewPingu failure %v", err)
	}
	defer pingu1.Close()

	if err := pingu1.RegisterWithRawAddr("127.0.0.1:9191"); err != nil {
		t.Fatalf("Register RegisterWithRawAddr failure : %v", err)
	}
	if err := pingu1.RegisterWithRawAddr("127.0.0.1:9192"); err != nil {
		t.Fatalf("Register RegisterWithRawAddr failure : %v", err)
	}
	if err := pingu1.RegisterWithRawAddr("127.0.0.1:9193"); err != nil {
		t.Fatalf("Register RegisterWithRawAddr failure : %v", err)
	}
	if len(pingu1.Pingus()) != 3 {
		t.Fatalf("PingPong failure got: %v, want : %v", len(pingu1.Pingus()), 3)
	}

	if err := pingu1.UnregisterWithRawAddr("127.0.0.1:9191"); err != nil {
		t.Fatalf("Register UnregisterWithRawAddr failure : %v", err)
	}
	if err := pingu1.UnregisterWithRawAddr("127.0.0.1:9192"); err != nil {
		t.Fatalf("Register UnregisterWithRawAddr failure : %v", err)
	}
	if err := pingu1.UnregisterWithRawAddr("127.0.0.1:9193"); err != nil {
		t.Fatalf("Register UnregisterWithRawAddr failure : %v", err)
	}
	if len(pingu1.Pingus()) != 0 {
		t.Fatalf("PingPong failure got: %v, want : %v", 0, len(pingu1.Pingus()))
	}
}

func TestBroadcastPingWithTicker(t *testing.T) {
	pingu1, err := pingu.NewPingu("127.0.0.1:9190", nil)
	if err != nil {
		t.Fatalf("BroadcastPingWithTicker NewPingu failure %v", err)
	}
	defer pingu1.Close()
	pingu2, err := pingu.NewPingu("127.0.0.1:9191", nil)
	if err != nil {
		t.Fatalf("BroadcastPingWithTicker NewPingu failure %v", err)
	}
	defer pingu2.Close()

	pingu1.Start()
	pingu2.Start()

	pingu1.RegisterWithRawAddr("127.0.0.1:9191")
	pingu1.RegisterWithRawAddr("127.0.0.1:9192")

	want := make(map[string]bool)
	want["127.0.0.1:9191"] = true
	want["127.0.0.1:9192"] = false

	cancel := pingu1.BroadcastPingWithTicker(*time.NewTicker(10 * time.Millisecond), 10*time.Millisecond)

	_ = cancel

	time.Sleep(50 * time.Millisecond)
	table := pingu1.PingTable()
	if !table["127.0.0.1:9191"] {
		t.Fatalf("BroadcastPingWithTicker invalid result: %v, want: %v", table, want)
	}
	if table["127.0.0.1:9192"] {
		t.Fatalf("BroadcastPingWithTicker invalid result: %v, want: %v", table, want)
	}

	close(cancel)
	time.Sleep(50 * time.Millisecond)
	table = pingu1.PingTable()
	if len(table) != 0 {
		t.Fatalf("BroadcastPingWithTicker invalid result length: %v, want: %v", len(table), 0)
	}
}
