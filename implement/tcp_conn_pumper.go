package implement

import (
	"bufio"
	"net"
	"network"
	"sync"
	"time"
)

const DefaultReadBuffSize = 8 * 1024
const DefaultWriteBuffSize = 16 * 1024
const MinMergedWriteBuffSize = 100 * 1024

type TcpConnPumper struct {
	net.Conn
	msgCh               chan []byte
	ctx                 *Context
	hub                 network.IHub
	LatestInterTime     time.Time //最新的一次交互时间
	readBuffSize        int
	mergedWriteBuffSize int
	disableMergedWrite  bool
	playerId            uint64
	playerMsgCh         chan []byte
}

func (c *TcpConnPumper) Close() {
	err := c.Conn.Close()
	if err != nil {

	}
}

func newTcpPumper(c net.Conn, ctx *Context, hub network.IHub) *TcpConnPumper {
	conn := &TcpConnPumper{
		Conn: c,
		ctx:  ctx,
		hub:  hub,
	}
	readBuffSize := DefaultReadBuffSize
	writeBuffSize := DefaultWriteBuffSize
	mergedWriteBuffSize := MinMergedWriteBuffSize

	if ctx.ReadBufferSize > 0 {
		readBuffSize = readBuffSize
	}
	if ctx.WriteBufferSize > 0 {
		writeBuffSize = ctx.WriteBufferSize
	}
	if ctx.MergedWriteBufferSize > mergedWriteBuffSize {
		mergedWriteBuffSize = ctx.MergedWriteBufferSize
	}
	c.(*net.TCPConn).SetReadBuffer(readBuffSize)
	c.(*net.TCPConn).SetWriteBuffer(writeBuffSize)
	return conn
}

func (c *TcpConnPumper) GetLatestInterTime() time.Time {
	return c.LatestInterTime
}

func (c *TcpConnPumper) Active() {
	if c.hub != nil {
		c.hub.ActiveConn(c)
	}
}

func (c *TcpConnPumper) Run() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		c.handleWrite()
		wg.Done()
	}()

	go func() {
		c.handleRead()
		wg.Done()
	}()
}

func (c *TcpConnPumper) Write(data []byte) (count int, err error) {
	return c.Conn.Write(data)
}

func (c *TcpConnPumper) handleWrite() {

	buff := NewDataBuff(c.mergedWriteBuffSize, !c.disableMergedWrite)
loop:
	for {
		select {
		case data := <-c.msgCh:
			c.SetWriteDeadline(time.Now().Add(WriteWait))
			rb, _ := buff.GetData(data, c.msgCh)
			_, err := c.Write(rb)
			if err != nil {
				break loop
			}
		}
	}
	c.Close()
}

func (c *TcpConnPumper) handleRead() {
	scanner := bufio.NewScanner(c.Conn)
	scanner.Buffer(make([]byte, c.readBuffSize), c.readBuffSize)
	scanner.Split(c.ctx.Splitter)
	for {
		if ok := scanner.Scan(); ok {
			data := scanner.Bytes()
			c.playerMsgCh <- data
		} else {
			break
		}
	}
}

func (c *TcpConnPumper) GetPlayerId() uint64 {
	return c.playerId
}
