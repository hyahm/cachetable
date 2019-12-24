package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		var msg string
		fmt.Print("what you say: ")
		fmt.Scanln(&msg)
		conn.Write([]byte(msg))
		var receive bytes.Buffer

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		receive.Write(buf[:n])
		fmt.Println("you say:", receive.String())
	}
	conn.Close()
}
