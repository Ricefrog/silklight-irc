package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"silklight/irc"
	"strings"
	"time"
)

const (
	libera     = "irc.libera.chat:6697"
	devdungeon = "irc.devdungeon.com:6667"
)

func main() {
	fmt.Println("Starting silklight-irc...")
	lainchan := irc.ServerInfo{"irc.lainchan.org", 6697}

	usingSSL := true
	var conn net.Conn
	var err error
	if usingSSL {
		conn, err = irc.ConnectSSL(lainchan)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn, err = irc.Connect(lainchan)
		if err != nil {
			log.Fatal(err)
		}
	}

	irc.Login(conn, "silklight")

	time.Sleep(2 * time.Second)
	irc.JoinChannel(conn, "#bots")
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Closing connection...")
		irc.SendMessage(conn, "#bots", "BOT: that's all folks")
		irc.Disconnect(conn)
	}()

	tp := textproto.NewReader(bufio.NewReader(conn))
	for {
		status, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(status, "PING") {
			irc.Pong(conn)
		}

		fmt.Println(status)
	}
}
