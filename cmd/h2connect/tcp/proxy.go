package tcp

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"

	"github.com/jasonmgh/playground/cmd/forfun/h2connect/h2"
	"golang.org/x/net/http2"
)

func forwardConn(conn net.Conn, h2t *http2.Transport, port string) {
	fwdConn := h2.StartConnectReq(context.Background(), h2t, "tcp-echo.fly.dev"+":"+port)
	if fwdConn == nil {
		log.Println("cannot connect")
		conn.Close()
		return
	}

	go io.Copy(fwdConn, conn)
	io.Copy(conn, fwdConn)
}

func StartProxy(host string, port string) {
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening [%s]...\n", host+":"+port)

	serverConn, err := net.Dial("tcp4", h2.Addr)
	if err != nil {
		log.Fatal(err)
	}

	h2t := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return serverConn, nil
		},
	}
	log.Println("connected to server host")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go forwardConn(conn, h2t, port)
	}
}
