package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for  {
		conn.Write([]byte("hello"))
		time.Sleep(time.Minute)
	}
}