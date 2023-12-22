package implement

import (
	"network"
	"time"
)

type BroadcastMessage struct {
	PlayerIds []uint64
	Data      []byte
}

type ConnChanData struct {
	Conn     network.IConn
	Category int
}

const (
	ChanCategoryInit = iota + 1
	ChanCategoryActive
	ChanCategoryClose
)

const (
	WriteWait      = 10 * time.Second
	PongWait       = 8 * time.Second //* 1000
	PingPeriod     = 5 * time.Second
	MaxMessageSize = 4096
)
