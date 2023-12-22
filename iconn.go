package network

import (
	"time"
)

type IConn interface {
	GetSession() ISession
	GetLatestInterTime() time.Time
	Close()
	Write([]byte) (count int, err error)
}
