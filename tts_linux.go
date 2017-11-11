// +build linux

package eintel

import (
  "fmt"
  "github.com/krig/go-sox"
)

type ttsLinux struct {
	messages chan IntelMessage
}

func NewTTS(messages chan IntelMessage) *ttsLinux {
	return &ttsLinux{
		messages: messages,
	}
}

func (v *ttsLinux) Run() {

  if !sox.Init() {
    panic("Failed to initialize SoX")
  }

  defer sox.Quit()

	for msg := range v.messages {
    err := PlayText(msg.Description)
    if err != nil {
      panic(fmt.Sprintf("%s \n\nPlease install vorbis-tools", err))
    }
	}
}

func PlayText(text string) (error) {
	content := speak(text)

  in := sox.OpenMemRead(content)
  defer in.Release()

	out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
	defer out.Release()

	chain := sox.CreateEffectsChain(in.Encoding(), out.Encoding())
  defer chain.Release()

	e := sox.CreateEffect(sox.FindEffect("input"))
	e.Options(in)
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()


	e = sox.CreateEffect(sox.FindEffect("output"))
	e.Options(out)
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()

  chain.Flow()
  return nil
}
