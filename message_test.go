package pingu_test

import (
	"bytes"
	"net"
	"net/netip"
	"reflect"
	"testing"

	"github.com/protocol-diver/pingu"
)

func TestParsePacket(t *testing.T) {
	type td struct {
		got    []byte
		expect pingu.Packet
		err    string
	}
	tdl := []td{
		{got: []byte{0, 2, 123, 125}, expect: packet(pingu.Ping, new(pingu.PingPacket)), err: ""},
		{got: []byte{1, 2, 123, 125}, expect: packet(pingu.Pong, new(pingu.PongPacket)), err: ""},
		{got: []byte{3, 2, 123, 125}, expect: nil, err: "invalid packet type: 3"},
		{got: []byte{4, 2, 123, 125}, expect: nil, err: "invalid packet type: 4"},
	}

	for _, td := range tdl {
		p, err := pingu.ParsePacket(td.got, nil)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("ParsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(p) != reflect.TypeOf(td.expect) {
			t.Fatalf("ParsePacket failure got: %v, want: %v", reflect.TypeOf(p), reflect.TypeOf(td.expect))
		}
	}

	uAddr := net.UDPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:1234"))
	for _, td := range tdl {
		p, err := pingu.ParsePacket(td.got, uAddr)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("ParsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(p) != reflect.TypeOf(td.expect) {
			t.Fatalf("ParsePacket failure got: %v, want: %v", reflect.TypeOf(p), reflect.TypeOf(td.expect))
		}
	}
}

func TestSuitablePack(t *testing.T) {
	//func SuitableUnpack(b []byte, packet Packet) error {
	type td struct {
		got    []byte
		expect pingu.Packet
		err    string
	}
	tdl := []td{
		{got: []byte{0, 2, 123, 125}, expect: packet(pingu.Ping, new(pingu.PingPacket)), err: ""},
		{got: []byte{1, 2, 123, 125}, expect: packet(pingu.Pong, new(pingu.PongPacket)), err: ""},
		{got: []byte{3, 2, 123, 125}, expect: nil, err: "invalid packet type: 3"},
		{got: []byte{4, 2, 123, 125}, expect: nil, err: "invalid packet type: 4"},
	}

	for _, td := range tdl {
		var pack pingu.Packet
		if _, ok := td.expect.(*pingu.PingPacket); ok {
			pack = new(pingu.PingPacket)
		} else {
			pack = new(pingu.PongPacket)
		}

		err := pingu.SuitablePack(td.got, pack)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("SuitableUnpack failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if reflect.TypeOf(pack) != reflect.TypeOf(td.expect) {
			t.Fatalf("SuitableUnpack failure got: %v, want: %v", reflect.TypeOf(pack), reflect.TypeOf(td.expect))
		}
	}
}

func TestSuitableUnpack(t *testing.T) {
	// SuitablePack(packet Packet) ([]byte, error) {
	type td struct {
		got    pingu.Packet
		expect []byte
		err    string
	}
	tdl := []td{
		{got: new(pingu.PingPacket), expect: []byte{0, 2, 123, 125}, err: ""},
		{got: new(pingu.PongPacket), expect: []byte{1, 2, 123, 125}, err: ""},
	}

	for _, td := range tdl {
		p, err := pingu.SuitableUnpack(td.got)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("ParsePacket failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
		if !bytes.Equal(p, td.expect) {
			t.Fatalf("SuitablePack failure got: %v, want: %v", p, td.expect)
		}
	}
}

func packet(typ byte, packet pingu.Packet) pingu.Packet {
	return packet
}
