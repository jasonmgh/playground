package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

// echo handler inspired by stdlib http2/http2_demo.go
type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ECHO", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "PUT required.", 400)
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		io.Copy(flushWriter{w}, r.Body)
	})

	l, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listening [%s]...\n", host)

	go func() {
		// only talk http2 directly, no need to use h2c for HTTP/1.1 compat
		server := http2.Server{}
		for {
			conn, _ := l.Accept()
			fmt.Println("handling new connection")

			go func() {
				server.ServeConn(conn, &http2.ServeConnOpts{Handler: mux})
				conn.Close()
			}()
		}
	}()
}
