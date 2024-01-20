package simulation

import "github.com/txthinking/socks5"

type SimulationServer struct {
	socks5.Server
	inner *socks5.Server
}

func NewSocks5Server(addr, ip, userName, password string, tcpTimeout, udpTimeout int) (*SimulationServer, error) {
	s5, err := socks5.NewClassicServer(addr, ip, userName, password, tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	x := &SimulationServer{
		inner: s5,
	}
	return x, nil
}

func (x *SimulationServer) ListenAndServe() error {
	handle := &SimulationHandle{}
	return x.inner.ListenAndServe(handle)
}

func (x *SimulationServer) Shutdown() error {
	return x.inner.Shutdown()
}
