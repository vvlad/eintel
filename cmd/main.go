package main

import (
	"github.com/ogier/pflag"
	"github.com/vvlad/eintel"
)

func main() {

	pflag.Parse()

	chat := eintel.NewChat("GOTG_Intel", "Private Chat (Yolla)")
	location := eintel.NewLocalWatcher(chat)
	parser := eintel.NewMessageParser(chat, location.Updates)
	voice := eintel.NewVoiceReport(parser.Messages)

	go chat.Run()
	go location.Run()
	go voice.Run()
	parser.Run()

}
