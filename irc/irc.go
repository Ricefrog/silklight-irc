package irc

import (
	"crypto/tls"
	"fmt"
	"net"
)

type ServerInfo struct {
	Domain string
	Port   int
}

func (s ServerInfo) String() string {
	return fmt.Sprintf("%s:%d", s.Domain, s.Port)
}

func ConnectSSL(server ServerInfo) (net.Conn, error) {
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", server.String(), conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Connect(server ServerInfo) (net.Conn, error) {
	conn, err := net.Dial("tcp", server.String())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Send a string through the socket connection
func SendString(conn net.Conn, msg string) {
	fmt.Fprintf(conn, "%s\r\n", msg)
}

func Login(conn net.Conn, nickname string) {
	SendString(conn, fmt.Sprintf("USER %s 0 * :silklight_irc", nickname))
	SendString(conn, fmt.Sprintf("NICK %s", nickname))
}

func JoinChannel(conn net.Conn, channelName string) {
	SendString(conn, fmt.Sprintf("JOIN %s", channelName))
}

func SendMessage(conn net.Conn, channelName, msg string) {
	SendString(conn, fmt.Sprintf("PRIVMSG %s %s", channelName, msg))
}

func Pong(conn net.Conn) {
	SendString(conn, "PONG")
}

func Disconnect(conn net.Conn) {
	SendString(conn, "QUIT adieu")
	conn.Close()
}
