package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

func startStreamReq(ctx context.Context, client *http.Client) {
	// write periodic timestamps to our request connection
	out, in := io.Pipe()

	// start an http2 request with our async time stream
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, io.NopCloser(out))
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("cannot connect, not ok: got [%d]\n", resp.StatusCode)
	}

	// start writing to our open request stream
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(in, "echo, time is: %v\n", time.Now())
		}
	}()

	// request body should be echo's of the input
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func startClient(ctx context.Context) {
	// open a new tcp connection to the server
	conn, err := net.Dial("tcp4", host)
	if err != nil {
		log.Fatal(err)
	}

	// create an http client, force http2, allow http connections without tls, reusing 1 tcp connection
	h2client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				// http2.Transport already implements connection pooling
				// but we're just proving a point here and re-using the conn directly
				return conn, nil
			},
		},
	}

	for c := 0; c < 3; c++ {
		go startStreamReq(ctx, h2client)
	}
}
