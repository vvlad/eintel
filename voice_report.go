package eintel

import (
	"os/exec"
)

type voiceReport struct {
	intelMessages chan IntelMessage
}

func NewVoiceReport(messages chan IntelMessage) *voiceReport {

	return &voiceReport{
		intelMessages: messages,
	}
}

func (v *voiceReport) Run() {

	for msg := range v.intelMessages {
		cmd := exec.Command("/usr/bin/say", msg.Description)
		cmd.Run()
		cmd.Wait()
	}

}
