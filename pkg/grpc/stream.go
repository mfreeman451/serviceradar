// pkg/grpc/stream.go

package grpc

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

// streamConn adapts a gRPC stream to net.Conn interface
type streamConn struct {
	grpc.ClientStream
}

func (s *streamConn) Close() error {
	// Implementation depends on your specific stream type
	// You'll need to use CloseSend to close the stream
	return s.CloseSend()
}

func (s *streamConn) Read(p []byte) (n int, err error) {
	// Implementation depends on your specific stream type
	// You'll need to use RecvMsg to receive data
	msg := &struct{ Data []byte }{}
	if err := s.RecvMsg(msg); err != nil {
		return 0, err
	}
	return copy(p, msg.Data), nil
}

func (s *streamConn) Write(p []byte) (n int, err error) {
	// Implementation depends on your specific stream type
	// You'll need to use SendMsg to send data
	if err := s.SendMsg(&struct{ Data []byte }{Data: p}); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *streamConn) LocalAddr() net.Addr                { return nil }
func (s *streamConn) RemoteAddr() net.Addr               { return nil }
func (s *streamConn) SetDeadline(t time.Time) error      { return nil }
func (s *streamConn) SetReadDeadline(t time.Time) error  { return nil }
func (s *streamConn) SetWriteDeadline(t time.Time) error { return nil }

// GrpcStreamToConn converts a gRPC stream to a net.Conn
func GrpcStreamToConn(stream grpc.ClientStream) net.Conn {
	return &streamConn{stream}
}
