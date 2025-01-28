// pkg/grpc/connection_context.go

package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc/peer"
)

// connAddr implements net.Addr for ClientConn.
type connAddr struct {
	conn *ClientConn
}

func (*connAddr) Network() string { return "grpc" }
func (*connAddr) String() string  { return "grpc-connection" }

// ConnectionFromContext extracts the ClientConn from the context.
func ConnectionFromContext(ctx context.Context) (*ClientConn, bool) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, false
	}

	// Since ClientConn now implements net.Addr, we can do this directly
	conn, ok := p.Addr.(*ClientConn)

	return conn, ok
}

// AsAddr converts a ClientConn to net.Addr for peer info.
func (c *ClientConn) AsAddr() net.Addr {
	return &connAddr{conn: c}
}
