// Copyright (c) 2022, Seungbae Yu <dbadoy4874@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pingu

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	ping = iota
	pong

	packetTypeIndex = 0
	packetSizeIndex = 1
	prefixSize      = 2
)

type packet interface {
	SetSender(s *net.UDPAddr)
	Sender() *net.UDPAddr
	Kind() byte
}

type pingPacket struct {
	sender *net.UDPAddr
}

type pongPacket struct {
	sender *net.UDPAddr
}

// parsePacket parses packets received by other pingus.
func parsePacket(d []byte, sender *net.UDPAddr) (packet, error) {
	var r packet
	switch d[packetTypeIndex] {
	case ping:
		r = new(pingPacket)
	case pong:
		r = new(pongPacket)
	default:
		return nil, fmt.Errorf("invalid packet type: %d", d[packetTypeIndex])
	}

	if err := suitablePack(d, r); err != nil {
		return nil, err
	}
	r.SetSender(sender)
	return r, nil
}

// suitablePack is the logic for parse the UDP Payload.
func suitablePack(b []byte, packet packet) error {
	if !isValidPacketType(b[packetTypeIndex]) {
		return fmt.Errorf("invalid packet type: %d", b[packetTypeIndex])
	}
	len := b[packetSizeIndex]
	byt := make([]byte, len)
	copy(byt[:], b[prefixSize:prefixSize+len])

	if err := json.Unmarshal(byt, packet); err != nil {
		return fmt.Errorf("invalid packet data: %v", err)
	}
	return nil
}

// suitableUnpack is change Packet to suitable protocol message.
// If send message, you must use this method.
func suitableUnpack(packet packet) ([]byte, error) {
	if !isValidPacketType(packet.Kind()) {
		return nil, fmt.Errorf("invalid packet type: %d", packet.Kind())
	}
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

func isValidPacketType(b byte) bool {
	switch b {
	case ping:
		return true
	case pong:
		return true
	default:
		return false
	}
}

func (p *pingPacket) SetSender(s *net.UDPAddr) { p.sender = s }
func (p *pingPacket) Sender() *net.UDPAddr     { return p.sender }
func (p *pingPacket) Kind() byte               { return ping }

func (p *pongPacket) SetSender(s *net.UDPAddr) { p.sender = s }
func (p *pongPacket) Sender() *net.UDPAddr     { return p.sender }
func (p *pongPacket) Kind() byte               { return pong }
