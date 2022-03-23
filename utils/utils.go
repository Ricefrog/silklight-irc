package utils

import (
	"fmt"
	"regexp"
	"silklight/irc"
	"silklight/styles"
	"strings"
	"time"
)

func RenderPRIVMSG(nick, msg string) string {
	tempStyle := styles.NickToStyle(nick)
	nick = fmt.Sprintf(
		"%s%s%s",
		styles.MagentaStyle.Render("<"),
		tempStyle.Render(nick),
		styles.MagentaStyle.Render(">"),
	)
	return fmt.Sprintf("%s %s", nick, msg)
}

func PrependTimestamp(str string) string {
	var b strings.Builder
	now := time.Now()
	fmt.Fprintf(&b, "[%s] %s %s", now.Format("15:04:05"), styles.Separator, str)
	return b.String()
}

var privmsg_exp = regexp.MustCompile(`:(.*)!.*PRIVMSG.*:(.*)`)

func IsPRIVMSG(raw string) bool {
	return privmsg_exp.MatchString(raw)
}

// :severian!severian@net-lcr.161.s088n6.IP PRIVMSG #bots :message
// Convert raw IRC PRIVMSG into something more palatable
func CleanPRIVMSG(raw string) string {
	matches := privmsg_exp.FindStringSubmatch(raw)
	return RenderPRIVMSG(matches[1], matches[2])
}

// removes prologye for messsages received when first joining server.
func CleanServerMsg() {
	fmt.Println("gay")
}

func CleanMessage(raw, clientName string, serverInfo irc.ServerInfo) string {
	if IsPRIVMSG(raw) {
		return CleanPRIVMSG(raw)
	}

	// clean up server messages
	reg := regexp.MustCompile(fmt.Sprintf(`:(%s).*%s(.*):(.*)`,
		serverInfo.Domain, clientName))
	if reg.MatchString(raw) {
		matches := reg.FindStringSubmatch(raw)
		if matches[2] == " " {
			return matches[3]
		} else {
			return matches[2] + ":" + matches[3]
		}
	}
	return raw
}

// TEST STUBS

func RegTestStub() {
	testMsg := ":severian!severian@net-lcr.161.s088n6.IP PRIVMSG #bots :message"
	fmt.Printf("%#v\n", privmsg_exp.MatchString(testMsg))
	fmt.Printf("%#v\n", privmsg_exp.MatchString("doesn't work"))
	fmt.Printf("%#v\n", privmsg_exp.FindStringSubmatch(testMsg))
}
