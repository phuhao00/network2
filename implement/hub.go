package implement

import (
	"network"
	"sync/atomic"
	"time"
)

type Hub struct {
	broadCastCh   chan *BroadcastMessage
	connCh        chan *ConnChanData
	idleTimeLimit time.Duration
	initConns     map[network.IConn]bool
	activeConns   map[uint64]network.IConn
	isClose       atomic.Bool
}

func NewHub(idleTimeLimit time.Duration) *Hub {
	h := &Hub{
		broadCastCh:   make(chan *BroadcastMessage, 100),
		connCh:        make(chan *ConnChanData, 1000),
		idleTimeLimit: idleTimeLimit,
		initConns:     make(map[network.IConn]bool),
		activeConns:   make(map[uint64]network.IConn),
	}
	return h
}

func (h *Hub) AddConn(conn network.IConn) {
	if h.isClose.Load() {
		return
	}
	if h.idleTimeLimit > 0 {
		h.pushChanData(conn, ChanCategoryInit)
	} else {
		h.pushChanData(conn, ChanCategoryActive)
	}
}

func (h *Hub) pushChanData(conn network.IConn, category int) {
	select {
	case h.connCh <- &ConnChanData{conn, category}:
	default:
	}
}

func (h *Hub) DelConn(conn network.IConn) {
	if h.isClose.Load() {
		return
	}
	h.pushChanData(conn, ChanCategoryClose)
}

func (h *Hub) ActiveConn(conn network.IConn) {
	if !h.isClose.Load() && h.idleTimeLimit > 0 {
		h.pushChanData(conn, ChanCategoryActive)
	}
}

func (h *Hub) Broadcast(sessionIds []uint64, data []byte) {
	h.broadCastCh <- &BroadcastMessage{PlayerIds: sessionIds, Data: data}
}

func (h *Hub) Stop() {
	if h.isClose.Load() {
		//todo
	} else {
		h.isClose.Store(true)
	}
}

func (h *Hub) Clear() {
	conns := make(map[network.IConn]bool)
	n := len(h.connCh)
	for i := 0; i < n; i++ {
		data := <-h.connCh
		if data.Category != ChanCategoryClose {
			conns[data.Conn] = true
		}
	}
	for conn := range h.initConns {
		conns[conn] = true
	}
	for _, conn := range h.activeConns {
		conns[conn] = true
	}
	for conn := range conns {
		conn.Close()
	}
}

func (h *Hub) GetActiveConnCount() int {
	return len(h.activeConns)
}

func (h *Hub) Run() {
	var ticker = time.NewTicker(3 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case data := <-h.connCh:
			conn := data.Conn
			switch data.Category {
			case ChanCategoryInit:
				h.initConns[conn] = true
			case ChanCategoryActive:
				if _, ok := h.initConns[conn]; ok {
					delete(h.initConns, conn)
				}
				h.activeConns[conn.GetPlayerId()] = conn
			case ChanCategoryClose:
				delete(h.initConns, conn)
				delete(h.activeConns, conn.GetPlayerId())
			}
		case message := <-h.broadCastCh:
			if len(message.PlayerIds) == 0 {
				for _, conn := range h.activeConns {
					conn.Write(message.Data)
				}
			} else {
				for _, id := range message.PlayerIds {
					conn := h.activeConns[id]
					if conn != nil {
						conn.Write(message.Data)
					}
				}
			}
		case <-ticker.C:
			if h.idleTimeLimit > 0 && len(h.initConns) > 0 {
				now := time.Now()
				for conn := range h.initConns {
					if now.Sub(conn.GetLatestInterTime()) > h.idleTimeLimit {
						delete(h.initConns, conn)
						conn.Close()
					}
				}
			}
		}
	}
}
