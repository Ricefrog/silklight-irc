package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os/exec"
	"silklight/frontend"
	futils "silklight/frontend/utils"
	"silklight/irc"
	"silklight/utils"
	"strings"
	"time"
)

var Quit bool = false

func main() {
	fmt.Println("Starting silklight-irc...")
	lainchan := irc.ServerInfo{"irc.lainchan.org", 6697}
	clientName := "silklight"

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

	m := &frontend.MainModel{Conn: conn, CurrentChannel: "#bots", ClientName: clientName}
	frontend.ClearScreen()
	p := frontend.Start(m)

	irc.Login(conn, clientName)

	time.Sleep(2 * time.Second)
	irc.JoinChannel(conn, "#bots")

	/*
		go func() {
			time.Sleep(120 * time.Second)
			fmt.Println("Closing connection...")
			irc.SendMessage(conn, "#bots", "BOT: that's all folks")
			p.Quit()
			irc.Disconnect(conn)
			os.Exit(0)
		}()
	*/

	tp := textproto.NewReader(bufio.NewReader(conn))
	for !futils.Quit {
		status, err := tp.ReadLine()
		if err != nil {
			conn.Close()
			log.Print("connection closed")
			p.Quit()
			break
		}

		if strings.HasPrefix(status, "PING") {
			irc.Pong(conn)
			continue
		}

		status = utils.CleanMessage(status, clientName, lainchan)

		p.Send(futils.AppendMsg(status + "\n"))
	}
	exec.Command("reset").Run()
	log.Print("silklight-irc quit.")
}
