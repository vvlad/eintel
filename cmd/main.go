package main

import (
	"github.com/ogier/pflag"
	"github.com/vvlad/eintel"
)

var (
	channels = []string{}
)

func init() {

}

func main() {

	pflag.Parse()

	chat := eintel.NewChat("GOTG_Intel", "Derzerek")
	location := eintel.NewLocalWatcher(chat)
	parser := eintel.NewMessageParser(chat, location.Updates)
	voice := eintel.NewLinuxTTS(parser.Messages)

	go chat.Run()
	go location.Run()
	go voice.Run()
	voice.PlayText("EIntel online")
	parser.Run()

}
