package h2

import "io"

type PipeConn struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (pc *PipeConn) Read(p []byte) (n int, err error) {
	return pc.reader.Read(p)
}

func (pc *PipeConn) Write(p []byte) (n int, err error) {
	return pc.writer.Write(p)

}

func (pc *PipeConn) Close() error {
	return pc.writer.Close()
}
