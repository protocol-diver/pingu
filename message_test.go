package pingu

import (
	"bytes"
	"net"
	"net/netip"
	"reflect"
	"testing"
)

func TestParsePacket(t *testing.T) {
	type td struct {
		got    []byte
		expect packet
		err    string
	}
	tdl := []td{
		{got: []byte{0, 2, 123, 125}, expect: new(pingPacket), err: ""},
		{got: []byte{1, 2, 123, 125}, expect: new(pongPacket), err: ""},
		{got: []byte{3, 2, 123, 125}, expect: nil, err: "invalid packet type: 3"},
		{got: []byte{4, 2, 123, 125}, expect: nil, err: "invalid packet type: 4"},
	}

	for _, td := range tdl {
		p, err := parsePacket(td.got, nil)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("parsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(p) != reflect.TypeOf(td.expect) {
			t.Fatalf("parsePacket failure got: %v, want: %v", reflect.TypeOf(p), reflect.TypeOf(td.expect))
		}
	}

	uAddr := net.UDPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:1234"))
	for _, td := range tdl {
		p, err := parsePacket(td.got, uAddr)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("parsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(p) != reflect.TypeOf(td.expect) {
			t.Fatalf("parsePacket failure got: %v, want: %v", reflect.TypeOf(p), reflect.TypeOf(td.expect))
		}
	}
}

func TestSuitablePack(t *testing.T) {
	//func SuitableUnpack(b []byte, packet Packet) error {
	type td struct {
		got    []byte
		expect packet
		err    string
	}
	tdl := []td{
		{got: []byte{0, 2, 123, 125}, expect: new(pingPacket), err: ""},
		{got: []byte{1, 2, 123, 125}, expect: new(pongPacket), err: ""},
		{got: []byte{3, 2, 123, 125}, expect: nil, err: "invalid packet type: 3"},
		{got: []byte{4, 2, 123, 125}, expect: nil, err: "invalid packet type: 4"},
	}

	for _, td := range tdl {
		var pack packet
		if _, ok := td.expect.(*pingPacket); ok {
			pack = new(pingPacket)
		} else {
			pack = new(pongPacket)
		}

		err := suitablePack(td.got, pack)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("suitableUnpack failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(pack) != reflect.TypeOf(td.expect) {
			t.Fatalf("suitableUnpack failure got: %v, want: %v", reflect.TypeOf(pack), reflect.TypeOf(td.expect))
		}
	}
}

func TestSuitableUnpack(t *testing.T) {
	// SuitablePack(packet Packet) ([]byte, error) {
	type td struct {
		got    packet
		expect []byte
		err    string
	}
	tdl := []td{
		{got: new(pingPacket), expect: []byte{0, 2, 123, 125}, err: ""},
		{got: new(pongPacket), expect: []byte{1, 2, 123, 125}, err: ""},
	}

	for _, td := range tdl {
		p, err := suitableUnpack(td.got)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("parsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if !bytes.Equal(p, td.expect) {
			t.Fatalf("suitablePack failure got: %v, want: %v", p, td.expect)
		}
	}
}
