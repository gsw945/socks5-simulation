package simulation

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

type SimulationHandle struct {
	socks5.DefaultHandle
	inner *socks5.DefaultHandle
}

// func (h *SimulationHandle) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
// 	return h.inner.TCPHandle(s, c, r)
// }

// ref: https://github.com/txthinking/socks5/blob/master/server.go#L323
// TCPHandle auto handle request. You may prefer to do yourself.
func (h *SimulationHandle) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		rc, err := r.Connect(c)
		if err != nil {
			return err
		}
		defer rc.Close()
		go func() {
			var bf [1024 * 2]byte
			for {
				if s.TCPTimeout != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
						return
					}
				}
				i, err := rc.Read(bf[:])
				if err != nil {
					return
				}
				if _, err := c.Write(bf[0:i]); err != nil {
					return
				}
			}
		}()
		var bf [1024 * 2]byte
		for {
			if s.TCPTimeout != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					return nil
				}
			}
			i, err := c.Read(bf[:])
			if err != nil {
				return nil
			}
			log.Printf("[TCPHandle]: Read() len=%d\n", i)
			n, err := rc.Write(bf[0:i])
			if err != nil {
				return nil
			}
			log.Printf("[TCPHandle]: Write() len=%d\n", n)
		}
		// return nil
	}
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(c, s.ServerAddr)
		if err != nil {
			return err
		}
		ch := make(chan byte)
		defer close(ch)
		s.AssociatedUDP.Set(caddr.String(), ch, -1)
		defer s.AssociatedUDP.Delete(caddr.String())
		io.Copy(io.Discard, c)
		if socks5.Debug {
			log.Printf("A tcp connection that udp %#v associated closed\n", caddr.String())
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}
