package network

type IPumper interface {
	WritePumper()
	ReadPumper()
}
