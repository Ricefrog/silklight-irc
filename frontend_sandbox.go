package main

import (
	"fmt"
	"silklight/frontend"
	futils "silklight/frontend/utils"
	"silklight/utils"
	"time"
)

func main() {
	m := &frontend.MainModel{}

	frontend.ClearScreen()
	p := frontend.Start(m)

	for i := 0; ; i++ {
		// SEND MESSAGE USING p.Send
		time.Sleep(time.Second)
		p.Send(futils.AppendMsg(utils.PrependTimestamp(fmt.Sprintf("Message %d\n", i))))
	}
}
