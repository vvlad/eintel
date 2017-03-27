package eintel

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type linuxTTS struct {
	messages chan IntelMessage
}

func NewLinuxTTS(messages chan IntelMessage) *linuxTTS {
	return &linuxTTS{
		messages: messages,
	}
}

func (v *linuxTTS) Run() {
	for msg := range v.messages {
		v.PlayText(msg.Description)
	}

}

func (v *linuxTTS) PlayText(text string) {
	content := speak(text)
	if len(content) > 0 {
		ioutil.WriteFile("/tmp/eintel-tts-pipe.ogg", content, 0400)
		cmd := exec.Command("ogg123", "/tmp/eintel-tts-pipe.ogg")
		cmd.Run()
		os.Remove("/tmp/eintel-tts-pipe.ogg")
	}
}
