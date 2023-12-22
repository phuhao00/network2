package implement

import (
	"net"
	"network"
	"runtime"
	"strings"
)

type Server struct {
	ctx *Context
	hub network.IHub
}

func NewServer(ctx *Context, hub network.IHub) *Server {
	s := &Server{
		ctx: ctx,
		hub: hub,
	}
	return s
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", s.ctx.ListenAddress)
	if err != nil {
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
			}
			break
		}
		conn := newTcpPumper(c, s.ctx, s.hub)
		conn.Run()
	}
}

func (s *Server) Stop() {
	s.hub.Stop()
}

func (s *Server) Broadcast(sessionIds []uint32, data []byte) {
	s.hub.Broadcast(sessionIds, data)
}
