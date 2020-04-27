package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var name = ""
var stop = make(chan int)
func main()  {
	serverAddr := "localhost:8888"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("dial err", err)
	}

	defer conn.Close()

	fmt.Printf("Make a nickname:")
	fmt.Scanf("%s", &name)
	fmt.Println("Welcome ", name)
	_, err = conn.Write([]byte("name|" + name))
	if err != nil {
		fmt.Println("write err", err)
	}
	go recv(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		ln, _, _ := reader.ReadLine()
		msg := string(ln)
		if msg == "quit" {

			_, err := conn.Write([]byte("quit|" + name))
			if err != nil {
				fmt.Println("write err", err)
			}
			stop<-0
			break
		}
		_, err := conn.Write([]byte("say|" + name + "|" + msg))
		if err != nil {
			fmt.Println("write err", err)
		}
	}
}

func recv(conn net.Conn)  {
	for {
		select {
		case <-stop:
			fmt.Println("Goodbye: ", name)
			return
		default:
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if n == 0 || err != nil {
				fmt.Println("read err", err)
			}
			strBuf := string(buffer[:n])
			fmt.Println(strBuf)
		}
	}
}