package eintel

import (
	"fmt"
)

type Message struct {
}

type messageParser struct {
	Messages         chan Message
	channelMesssages chan ChannelMessage
}

func NewMessageParser(channelMesssages chan ChannelMessage) *messageParser {
	return &messageParser{
		channelMesssages: channelMesssages,
		Messages:         make(chan Message),
	}
}

func (m *messageParser) Run() {
	for channelMessage := range m.channelMesssages {
		m.process(channelMessage)
	}
}

func (m *messageParser) process(msg ChannelMessage) {
	fmt.Printf("#%s %s\n", msg.Channel, msg.Message)
	m.Messages <- Message{}
}
