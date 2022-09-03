package pingu

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const (
	Ping = iota
	Pong

	packetTypeIndex = 0
	packetSizeIndex = 1
	prefixSize      = 2
)

type Packet interface {
	SetSender(s *net.UDPAddr)
	Sender() *net.UDPAddr
	SetKind(k byte)
	Kind() byte
}

type PingPacket struct {
	sender *net.UDPAddr
	kind   byte
}

type PongPacket struct {
	sender *net.UDPAddr
	kind   byte
}

func ParsePacket(d []byte, sender *net.UDPAddr) (Packet, error) {
	var r Packet
	switch d[packetTypeIndex] {
	case Ping:
		r = new(PingPacket)
	case Pong:
		r = new(PongPacket)
	default:
		return nil, fmt.Errorf("invalid packet type: %d", d[packetTypeIndex])
	}

	if err := SuitableUnpack(d, r); err != nil {
		return nil, err
	}
	r.SetSender(sender)
	r.SetKind(d[packetTypeIndex])
	return r, nil
}

// SuitablePack is change Packet to suitable protocol message.
// If send message, you must use this method.
func SuitablePack(packet Packet) ([]byte, error) {
	b, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}
	result := make([]byte, len(b)+prefixSize)
	result[packetTypeIndex] = packet.Kind()
	result[packetSizeIndex] = byte(len(b))
	copy(result[2:], b[:])
	return result, nil
}

// This is the logic for parse the Payload
func SuitableUnpack(b []byte, packet Packet) error {
	len := b[packetSizeIndex]
	byt := make([]byte, len)
	copy(byt[:], b[prefixSize:prefixSize+len])
	if err := json.Unmarshal(byt, packet); err != nil {
		return errors.New("invalid packet data")
	}
	return nil
}

func (p *PingPacket) SetSender(s *net.UDPAddr) { p.sender = s }
func (p *PingPacket) Sender() *net.UDPAddr     { return p.sender }
func (p *PingPacket) SetKind(k byte)           { p.kind = k }
func (p *PingPacket) Kind() byte               { return p.kind }

func (p *PongPacket) SetSender(s *net.UDPAddr) { p.sender = s }
func (p *PongPacket) Sender() *net.UDPAddr     { return p.sender }
func (p *PongPacket) SetKind(k byte)           { p.kind = k }
func (p *PongPacket) Kind() byte               { return p.kind }