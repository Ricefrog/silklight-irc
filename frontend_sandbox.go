package main

import (
	"fmt"
	"log"
	"os"
	"silklight/frontend"
	futils "silklight/frontend/utils"
	"time"
)

func main() {
	filename := "deleteme"
	os.Remove("deleteme")
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := &frontend.MainModel{}

	frontend.ClearScreen()
	p := frontend.Start(m)

	message := "$"
	for i := 0; ; i++ {
		// SEND MESSAGE USING p.Send
		time.Sleep(time.Second / 5.0)
		p.Send(futils.AppendMsg(message))
		p.Send(futils.AppendMsgToTextBox(fmt.Sprintf("Msg length: %d", len(message))))
		message += "$"
		file.WriteString(fmt.Sprintf("lines: %d\n", m.ViewPort.NumLines()))
		file.WriteString(m.ViewPort.PrintLiterals())
		file.WriteString(m.ViewPort.PrintLines())
	}
}
