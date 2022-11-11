package main

import (
	"github.com/jasonmgh/playground/cmd/forfun/h2connect/h2"
	"github.com/jasonmgh/playground/cmd/forfun/h2connect/tcp"
)

// Demo of proxying multiple tcp connections to a target server over http2 frames.
// Thought being to multiplex client tcp connections from the proxy to the app server
// over one open socket / tcp connection.
// Probably has too many performance tradeoffs, but what the heck.
// |------- PROXY SERVER -------|  <--->  |-------------- APP SERVER --------------|
// client 1 --\                                          /--- tcp-echo.fly.dev:5001
// .           :5001 --> h2_client ---> :443 h2_server -->
// client 2 --/                                          \---- tcp-echo.fly.dev:5001
// .
// client 3 --\                                          /--- tcp-echo.fly.dev:5002
// .           :5002 --> h2_client ---> :443 h2_server -->
// client 4 --/                                          \---- tcp-echo.fly.dev:5002
func main() {
	h2.StartServer()

	go tcp.StartProxy("localhost", "5001")
	tcp.StartProxy("localhost", "5002")

	// from shell, run nc localhost 5001, should get fly echo response.
}
