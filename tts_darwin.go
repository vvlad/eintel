// +build darwin

package eintel

import (
//	"fmt"
	"os/exec"
)

type ttsDarwin struct {
	messages chan IntelMessage
}

func NewTTS(messages chan IntelMessage) *ttsDarwin {
	return &ttsDarwin{
		messages: messages,
	}
}

func (v *ttsDarwin) Run() {
	for msg := range v.messages {
		v.PlayText(msg.Description)
	}

}

func (v *ttsDarwin) PlayText(text string) {
	if len(text) > 0 {
		cmd := exec.Command("say", "-r","280",text)
		cmd.Run()
	}
}
