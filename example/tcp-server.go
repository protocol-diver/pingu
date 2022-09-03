package example

import (
	"net"
	"net/netip"
	"time"

	"github.com/dbadoy/pingu"
)

type TCPServer struct {
	conn  *net.TCPConn
	pingu *pingu.Pingu
}

func NewServer(conn *net.TCPConn) *TCPServer {
	tempPingu, err := pingu.NewPingu(pingu.DefaultAddress(), nil)
	if err != nil {
		return nil
	}
	return &TCPServer{
		conn:  conn,
		pingu: tempPingu,
	}
}

func (s *TCPServer) heartbeatLoop(addrs []string, ticker *time.Ticker) (chan struct{}, error) {
	s.pingu.Start()
	for _, addr := range addrs {
		rawAddr := netip.MustParseAddrPort(addr)
		if err := s.pingu.Register(rawAddr.String()); err != nil {
			return nil, err
		}
	}
	return s.pingu.BroadcastPingWithTicker(*ticker, 3*time.Second), nil
}

func (s *TCPServer) stat() map[string]bool {
	return s.pingu.PingTable()
}

// ... TCP controll logics ~
