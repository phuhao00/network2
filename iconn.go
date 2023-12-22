package network

import (
	"time"
)

type IConn interface {
	GetLatestInterTime() time.Time
	Close()
	Write([]byte) (count int, err error)
	GetPlayerId() uint64
}
