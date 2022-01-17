package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

const (
	lainchan   = "irc.lainchan.org:6697"
	libera     = "irc.libera.chat:6697"
	devdungeon = "irc.devdungeon.com:6667"
)

func main() {
	fmt.Println("Starting silklight-irc...")

	usingSSL := true
	var conn net.Conn
	var err error
	if usingSSL {
		conf := &tls.Config{}
		conn, err = tls.Dial("tcp", lainchan, conf)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn, err = net.Dial("tcp", lainchan)
		if err != nil {
			log.Fatal(err)
		}
	}

	tp := textproto.NewReader(bufio.NewReader(conn))
	for {
		status, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(status)
	}

	conn.Close()
}
