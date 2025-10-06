package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from ", conn.RemoteAddr().String())

	var sizeBuf [4]byte
	_, err := io.ReadFull(conn, sizeBuf[:])
	if err != nil {
		fmt.Println("Error reading size: ", err)
		return
	}
	size := binary.BigEndian.Uint32(sizeBuf[:])

	var errorCode int16 = 0
	var corrID uint32
	if size > 0 {
		if size < 8 {
			fmt.Println("Request too short for header v2: ", size)
			_, _ = io.CopyN(io.Discard, conn, int64(size))
		} else {
			var headerPrefix [8]byte
			_, err := io.ReadFull(conn, headerPrefix[:])
			if err != nil {
				fmt.Println("Error reading header: ", err)
				return
			}

			apiVersion := binary.BigEndian.Uint16(headerPrefix[2:4])
			corrID = binary.BigEndian.Uint32(headerPrefix[4:8])

			if apiVersion > 4 {
				errorCode = 35 // Unsupported Version
			}

			remaining := int(size) - 8
			if remaining > 0 {
				_, err := io.CopyN(io.Discard, conn, int64(remaining))
				if err != nil {
					fmt.Println("Error draining remaining payload: ", err)
					return
				}
			}
		}
	}

	resp := make([]byte, 10)
	binary.BigEndian.PutUint32(resp[0:4], 0)
	binary.BigEndian.PutUint32(resp[4:8], corrID)
	binary.BigEndian.PutUint16(resp[8:10], uint16(errorCode))
	_, err = conn.Write(resp)
	if err != nil {
		fmt.Println("Error writing response: ", err)
		return
	}
	return
}
