package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	// 定义报文头部固定字段长度
	StartBytesLen  = 2
	LengthBytesLen = 2
	AFNBytesLen    = 1
	DeviceIDLen    = 10
	SocketNumLen   = 1
	HeaderLen      = StartBytesLen + LengthBytesLen + AFNBytesLen + DeviceIDLen + SocketNumLen
)

func main() {
	// 启动 TCP Server 监听在本地的8080端口
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("TCP Server started and listening on :8080")

	for {
		// 接受客户端连接
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		// 启动一个 goroutine 处理客户端连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 定义缓冲区和数据缓冲区
	for {
		// 从连接中读取数据到缓冲区
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			return
		}
		//
		//fmt.Printf("buffer Bytes: %X\n, n %d\n", buffer, n)
		// 将读取到的数据写入数据缓冲区
		var dataBuffer bytes.Buffer
		dataBuffer.Write(buffer[:n])
		fmt.Printf("dataBuffer length %d\n", dataBuffer.Len())

		// 处理数据缓冲区中的数据
		for {
			// 检查数据缓冲区中是否有足够的数据来解析报文头部
			if dataBuffer.Len() < HeaderLen {
				break
			}

			// 检查起始域是否有效
			startBytes := dataBuffer.Bytes()[:StartBytesLen]
			if !bytes.Equal(startBytes, []byte{0xaa, 0xaa}) {
				fmt.Println("Invalid start bytes")
				// 如果起始域无效，跳过一个字节继续检查
				dataBuffer.Next(1)
				continue
			}

			// 读取数据长度字段
			lengthBytes := dataBuffer.Bytes()[StartBytesLen : StartBytesLen+LengthBytesLen]
			length := binary.BigEndian.Uint16(lengthBytes)

			// 计算完整报文的总长度
			totalLen := int(StartBytesLen + LengthBytesLen + length)
			fmt.Printf("dataBuffer length %d, totalLen %d", dataBuffer.Len(), totalLen)
			// 检查是否有足够的数据来解析完整报文
			if dataBuffer.Len() < totalLen {
				return
			}

			// 提取完整报文
			message := dataBuffer.Next(totalLen)
			// 打印原始消息的十六进制表示
			fmt.Printf("Raw Message: %X\n", message)
			fmt.Println("----- Parsed Message -----")
			// 打印解析后的消息
			printMessage(message)
			fmt.Println("--------------------------")
		}
	}
}

func printMessage(message []byte) {
	// 提取并打印报文的各个部分
	startBytes := message[:StartBytesLen]
	lengthBytes := message[StartBytesLen : StartBytesLen+LengthBytesLen]
	length := binary.BigEndian.Uint16(lengthBytes)
	afn := message[StartBytesLen+LengthBytesLen]
	deviceID := message[StartBytesLen+LengthBytesLen+AFNBytesLen : StartBytesLen+LengthBytesLen+AFNBytesLen+DeviceIDLen]
	socketNum := message[StartBytesLen+LengthBytesLen+AFNBytesLen+DeviceIDLen]
	dataUnit := message[HeaderLen : HeaderLen+length]

	fmt.Printf("Start Bytes: %X\n", startBytes)
	fmt.Printf("Length: %d\n", length)
	fmt.Printf("AFN: %X\n", afn)
	fmt.Printf("Device ID: %s\n", string(deviceID))
	fmt.Printf("Socket Number: %d\n", socketNum)
	fmt.Printf("Data Unit: %X\n", dataUnit)
}
