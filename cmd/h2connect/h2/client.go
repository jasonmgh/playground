package h2

import (
	"context"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func StartConnectReq(ctx context.Context, h2t *http2.Transport, authority string) io.ReadWriteCloser {
	pr, pw := io.Pipe()

	// start an http2 request with our async time stream
	req, err := http.NewRequestWithContext(ctx, http.MethodConnect, "https://"+authority, io.NopCloser(pr))
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := h2t.RoundTrip(req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("finished resp")

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("cannot connect, not ok: got [%d]\n", resp.StatusCode)
	}

	return &PipeConn{
		reader: resp.Body,
		writer: pw,
	}
}
