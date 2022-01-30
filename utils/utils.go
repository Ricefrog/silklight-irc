package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func PrependTimestamp(str string) string {
	var b strings.Builder
	now := time.Now()
	fmt.Fprintf(&b, "%s %s", now.Format("15:04:05"), str)
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
	ret := fmt.Sprintf("<%s> %s", matches[1], matches[2])
	return ret
}

func RegTestStub() {
	testMsg := ":severian!severian@net-lcr.161.s088n6.IP PRIVMSG #bots :message"
	fmt.Printf("%#v\n", privmsg_exp.MatchString(testMsg))
	fmt.Printf("%#v\n", privmsg_exp.MatchString("doesn't work"))
	fmt.Printf("%#v\n", privmsg_exp.FindStringSubmatch(testMsg))
}
