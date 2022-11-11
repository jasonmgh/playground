package main

import (
	"context"
	"fmt"
	"time"
)

var (
	host = "localhost:8081"
	url  = fmt.Sprintf("http://%s/ECHO", host)
)

// h2stream.go is a demo of multiple full-duplex streaming requests being made
// to a http server over 1 tcp connection.
// Kind of a precursor to cmd/h2connect -- this app is just a POC of h2 streaming
// w/o the listening tcp proxy.
func main() {
	startServer()

	// run the request loop for 3s and then exit the program
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	startClient(ctx)

	<-ctx.Done()
}
