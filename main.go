package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"github.com/urfave/cli/v3"

	"socks5-simulation/simulation"
)

func httpServer(g *runnergroup.RunnerGroup) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		simulation.EnableDelay = !simulation.EnableDelay
		simulation.EnableLoss = !simulation.EnableLoss
		fmt.Fprintf(w, "Hello World!")
	})
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: http.DefaultServeMux,
	}
	// Create a context with a 5-second timeout for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	g.Add(&runnergroup.Runner{
		Start: func() error {
			return srv.ListenAndServe()
		},
		Stop: func() error {
			// Shutdown the HTTP server gracefully.
			return srv.Shutdown(ctx)
		},
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	g := runnergroup.New()
	cmd := cli.Command{
		Name:                  "socks5-simulation",
		Version:               "v0.0.1",
		Usage:                 "socks5-simulation",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log",
				Usage: "Enable log. A valid value is file path or 'console'. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true",
			},
			&cli.StringFlag{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "Socks5 server listen address, like: :1080 or 1.2.3.4:1080",
			},
			&cli.StringFlag{
				Name:  "username",
				Usage: "User name, optional",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Password, optional",
			},
			&cli.BoolFlag{
				Name:  "limitUDP",
				Usage: "The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.",
			},
			&cli.IntFlag{
				Name:  "tcpTimeout",
				Value: 0,
				Usage: "Connection deadline time (s)",
			},
			&cli.IntFlag{
				Name:  "udpTimeout",
				Value: 60,
				Usage: "Connection deadline time (s)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			listen := cmd.String("listen")
			if listen == "" {
				return cli.ShowSubcommandHelp(cmd)
			}
			host, _, err := net.SplitHostPort(listen)
			if err != nil {
				return err
			}
			var ip string
			if host != "" {
				ip = host
			} else {
				ip = "127.0.0.1"
			}
			username := cmd.String("username")
			password := cmd.String("password")
			tcpTimeout := int(cmd.Int("tcpTimeout"))
			udpTimeout := int(cmd.Int("udpTimeout"))
			s5, err := simulation.NewSocks5Server(listen, ip, username, password, tcpTimeout, udpTimeout)
			if err != nil {
				return err
			}
			s5.LimitUDP = cmd.Bool("limitUDP")
			g.Add(&runnergroup.Runner{
				Start: func() error {
					return s5.ListenAndServe()
				},
				Stop: func() error {
					return s5.Shutdown()
				},
			})
			return nil
		},
		ShellComplete: func(ctx context.Context, cmd *cli.Command) {
			l := cmd.VisibleFlags()
			for _, v := range l {
				fmt.Println("--" + v.Names()[0])
			}
		},
	}
	if os.Getenv("SOCKS5_DEBUG") != "" {
		socks5.Debug = true
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
	// http server
	httpServer(g)
	if len(g.Runners) == 0 {
		return
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		g.Done()
	}()
	if err := g.Wait(); err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
}
