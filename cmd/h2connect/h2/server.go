package h2

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

const Addr = "localhost:8443"

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

func forwardCONNECTor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		log.Println("not connect")
		http.Error(w, "CONNECT required", 400)
		return
	}

	log.Printf("starting conn to %s\n", r.Host)
	conn, err := net.Dial("tcp4", r.Host)
	if err != nil {
		log.Fatal(err)
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	go io.Copy(flushWriter{w}, conn)
	io.Copy(flushWriter{conn}, r.Body)
}

func StartServer() {
	l, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listening [%s]...\n", Addr)

	go func() {
		// only talk http2 directly, no need to use h2c for HTTP/1.1 compat
		server := http2.Server{}
		for {
			conn, _ := l.Accept()
			fmt.Println("handling new connection")

			go func() {
				server.ServeConn(conn, &http2.ServeConnOpts{
					Handler: http.HandlerFunc(forwardCONNECTor),
				})
				conn.Close()
			}()
		}
	}()
}
