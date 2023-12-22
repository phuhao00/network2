package implement

import "bytes"

type DataBuff struct {
	buff     *bytes.Buffer
	buffSize int
	enabled  bool
}

func NewDataBuff(buffSize int, enabled bool) *DataBuff {
	if !enabled {
		return &DataBuff{enabled: enabled}
	}
	return &DataBuff{
		buff:     bytes.NewBuffer(make([]byte, buffSize+2*1024)),
		buffSize: buffSize,
	}
}

func (this *DataBuff) GetData(data []byte, c <-chan []byte) ([]byte, int) {
	if !this.enabled || len(c) == 0 {
		return data, 1
	}

	buff, buffSize, count := this.buff, this.buffSize, 0
	buff.Reset()
	for {
		data = <-c
		count++
		buff.Write(data)
		if len(c) == 0 || buff.Len() >= buffSize {
			break
		}
	}
	return buff.Bytes(), count
}
