package main

import (
	"fmt"
	"net"
	"strings"
)

var ClientMap = make(map[string]net.Conn)

func main()  {
	fmt.Println("Server started")
	l, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Started error", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error", err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	for  {
		buf := make([]byte, 1024)
		lenRev, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", conn.RemoteAddr(), " : ",err)
		}
		strBuf := string(buf[:lenRev])
		//fmt.Println("Message: ", strBuf)
		msg := strings.Split(strBuf, "|")
		switch msg[0] {
		case "name":
			fmt.Println(conn.RemoteAddr(), "-->", msg[1], " connected")
			for k, v := range ClientMap {
				if k != msg[1] {
					_, err := v.Write([]byte(msg[1] + ": joined"))
					if err != nil {
						fmt.Println("write to ", k, "err", err)
					}
				}
			}
			ClientMap[msg[1]] = conn
		case "say":
			fmt.Println("Got message from ", msg[1])
			for k, v := range ClientMap {
				if k != msg[1] {
					_, err := v.Write([]byte(msg[1] + ": " + msg[2]))
					if err != nil {
						fmt.Println("write to ", k, "err", err)
					}
				}
			}
		case "quit":
			fmt.Println(msg[1], " has quited")
			for k, v := range ClientMap {
				if k != msg[1] {
					_, err := v.Write([]byte(msg[1] + ": quited"))
					if err != nil {
						fmt.Println("write to ", k, "err", err)
					}
				}
			}
			delete(ClientMap, msg[1])
			conn.Close()
			return
		}
	}
}

