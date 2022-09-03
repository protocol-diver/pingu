package example

import (
	"net"
	"net/netip"
	"pingu"
	"time"
)

type TCPServer struct {
	conn  *net.TCPConn
	pingu Pingu
}

func NewServer(conn *net.TCPConn) *TCPServer {
	uconn, err := net.DialUDP("udp", nil, nil)
	if err != nil {
		return nil
	}
	tempPingu := pingu.NewPingu(uconn, pingu.Config{})
	return &TCPServer{
		conn:  conn,
		pingu: tempPingu,
	}
}

func (s *TCPServer) heartbeatLoop(addrs []string, ticker *time.Ticker) (chan struct{}, error) {
	s.pingu.Start()
	for _, addr := range addrs {
		rawAddr := netip.MustParseAddrPort(addr)
		if err := s.pingu.Register(net.UDPAddrFromAddrPort(rawAddr)); err != nil {
			return nil, err
		}
	}
	return s.pingu.BroadcastPingWithTicker(*ticker, 3*time.Second)
}

func (s *TCPServer) stat() map[string]bool {
	return s.pingu.PingTable()
}

// ... TCP controll logics ~
