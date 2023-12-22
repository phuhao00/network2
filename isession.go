package network

type ISession interface {
	GetId() uint32
	GetConn() IConn
	OnConnect(conn IConn)
	OnClose(conn IConn)
	OnReceive(conn IConn, data []byte)
}
