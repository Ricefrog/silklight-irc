package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"silklight/frontend"
	futils "silklight/frontend/utils"
	"silklight/irc"
	"silklight/utils"
	"strings"
	"time"
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

	m := &frontend.MainModel{Conn: conn, CurrentChannel: "#bots", ClientName: "silklight"}
	frontend.ClearScreen()
	p := frontend.Start(m)

	irc.Login(conn, "silklight")

	time.Sleep(2 * time.Second)
	irc.JoinChannel(conn, "#bots")

	quit := false
	go func() {
		time.Sleep(120 * time.Second)
		fmt.Println("Closing connection...")
		irc.SendMessage(conn, "#bots", "BOT: that's all folks")
		p.Quit()
		irc.Disconnect(conn)
		os.Exit(0)
	}()

	tp := textproto.NewReader(bufio.NewReader(conn))
	for !quit {
		status, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(status, "PING") {
			irc.Pong(conn)
			continue
		}

		if utils.IsPRIVMSG(status) {
			status = utils.CleanPRIVMSG(status)
		}

		p.Send(futils.AppendMsg(status + "\n"))
	}
}
