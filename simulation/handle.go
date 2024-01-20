package simulation

import (
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

var delayMin = 1   // min delay, unit is ms, >= 1
var delayMax = 300 // max delay, unit is ms, >= delayMin

type SimulationHandle struct {
	socks5.DefaultHandle
	inner *socks5.DefaultHandle
}

// func (h *SimulationHandle) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
// 	return h.inner.TCPHandle(s, c, r)
// }

func (h *SimulationHandle) randSleep(bound string) {
	// Generate a random delay between 20 and 300 milliseconds.
	min := math.Max(1, float64(delayMin))
	max := math.Max(min, float64(delayMax)-min)
	delay := int(min) + rand.Intn(int(max))
	log.Printf("[randSleep]: [%s] delay=%dms\n", bound, delay)
	// Send the current time to the channel after the delay.
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ref: https://github.com/txthinking/socks5/blob/master/server.go#L323
// TCPHandle auto handle request. You may prefer to do yourself.
func (h *SimulationHandle) TCPHandle(s *socks5.Server, upstream *net.TCPConn, r *socks5.Request) error {
	// Create a random number generator.
	rand.Seed(time.Now().UnixNano())
	if r.Cmd == socks5.CmdConnect {
		client, err := r.Connect(upstream)
		if err != nil {
			return err
		}
		defer client.Close()
		// inbound: read from client and write to upstream
		go func() {
			var bf [1024 * 2]byte
			for {
				if s.TCPTimeout != 0 {
					if err := client.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
						return
					}
				}
				i, err := client.Read(bf[:])
				if err != nil {
					return
				}
				log.Printf("[TCPHandle]: inbound Read(): len=%d\n", i)
				h.randSleep("inbound") // simulate network inbound delay
				n, err := upstream.Write(bf[0:i])
				if err != nil {
					return
				}
				log.Printf("[TCPHandle]: inbound Write(): len=%d\n", n)
			}
		}()
		// outbound: read from upstream and write to client
		var bf [1024 * 2]byte
		for {
			if s.TCPTimeout != 0 {
				if err := upstream.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					return nil
				}
			}
			i, err := upstream.Read(bf[:])
			if err != nil {
				return nil
			}
			log.Printf("[TCPHandle]: outbound Read() len=%d\n", i)
			h.randSleep("outbound") // simulate network outbound delay
			n, err := client.Write(bf[0:i])
			if err != nil {
				return nil
			}
			log.Printf("[TCPHandle]: outbound Write() len=%d\n", n)
		}
		// return nil
	}
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(upstream, s.ServerAddr)
		if err != nil {
			return err
		}
		ch := make(chan byte)
		defer close(ch)
		s.AssociatedUDP.Set(caddr.String(), ch, -1)
		defer s.AssociatedUDP.Delete(caddr.String())
		io.Copy(io.Discard, upstream)
		if socks5.Debug {
			log.Printf("A tcp connection that udp %#v associated closed\n", caddr.String())
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}
