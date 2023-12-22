package network

type IHub interface {
	Run()
	AddConn(conn IConn)
	DelConn(conn IConn)
	ActiveConn(conn IConn)
	Broadcast(sessionIds []uint32, data []byte)
	Stop()
	Clear()
	GetActiveConnCount() int
}
