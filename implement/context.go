package implement

import (
	"bufio"
	"time"
)

type Context struct {
	Splitter              bufio.SplitFunc
	IPChecker             func(ip string) bool
	IdleTimeAfterOpen     time.Duration
	ReadBufferSize        int
	WriteBufferSize       int
	UseNoneBlockingChan   bool
	ChanSize              int
	MaxMessageSize        int
	MergedWriteBufferSize int
	DisableMergedWrite    bool
	EnableStatistics      bool
	Extra                 interface{}
	ListenAddress         string
}
