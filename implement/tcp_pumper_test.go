package implement

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestTcpPumperClient(t *testing.T) {
	// 创建 TCP 地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("无法解析 TCP 地址:", err)
		os.Exit(1)
	}

	// 建立 TCP 连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("无法连接到服务器:", err)
		os.Exit(1)
	}
	defer conn.Close()
	for {
		// 向服务器发送数据
		message := "Hello, server!"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("发送数据失败:", err)
			os.Exit(1)
		}

		// 从服务器接收响应
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("接收数据失败:", err)
			os.Exit(1)
		}

		// 处理响应
		response := string(buffer[:bytesRead])
		fmt.Println("收到服务器响应:", response)
	}
}

func TestTcpPumperServer(t *testing.T) {

}
