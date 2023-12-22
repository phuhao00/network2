package network

type IMessageChan interface {
	GetInCh() chan<- []byte
	GetOutCh() <-chan []byte
	Len() int
}
