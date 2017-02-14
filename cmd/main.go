package main

import (
	"fmt"
	"github.com/vvlad/eintel"
	"github.com/vvlad/eintel/universe"
)

func main() {

	universe.Load()

	chat := eintel.NewChatLogsWatcher("GOTG_Intel")
	parser := eintel.NewMessageParser(chat.Messages)

	go chat.Run()
	parser.Run()

	for message := range parser.Messages {
		fmt.Printf("%v\n", message)
	}

}
