package eintel

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/krig/go-sox"
	"io/ioutil"
  "strings"
)

type voice struct {
  threats chan ThreatMessage
}

func NewVoiceNotification(messages chan ThreatMessage) *voice {
  v := &voice{
    threats: messages,
  }

  go v.run()
  return v
}

func (v *voice)run() {
	if !sox.Init() {
		panic("Failed to initialize SoX")
	}

  defer sox.Quit()
  playText("EIntel Online")

  for message := range v.threats {
    if message.Level > ThreatLevelIrelevant {
      text := v.formatMessage(message)
      playText(text)
    }else{
      log.Debugf("[%s] discarded message %s because it's thread level is %d", message.PlayerName, message.Line, message.Level)
    }

    log.Info("Notification sent")
  }

}

func (v *voice) formatMessage(message ThreatMessage) string {
  if message.Level == ThreatLevelCleared {
    return fmt.Sprintf("%s is clear", formatLocation(message.RelatedSystem.Name))
  } else if message.Jumps == 0 {
    return "Threat in local. Dock Dock Dock"
  } else if message.Jumps == 1 {
    return fmt.Sprintf("Threat one jump away from %s", message.PlayerName)
  } else if message.Jumps > 1 {
    return fmt.Sprintf("Threat %d jumps away from %s", message.Jumps, message.PlayerName)
  }
  return "Thread message"
}

func playText(text string) error {
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

var (
	emptyBytes = make([]byte, 0)
)

func speak(text string) []byte {
	cacheKey := fmt.Sprintf("aws:poly:%s", MD5Hash(text))
	if data, inCache := cache.Get(cacheKey); inCache {
		return data
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("eu-west-1")},
	}))
	svc := polly.New(sess)
	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("ogg_vorbis"), // Required
		TextType:     aws.String("ssml"),
		Text:         aws.String(fmt.Sprintf("<speak><prosody rate=\"x-fast\">%s</prosody></speak>", text)), // Required
		VoiceId:      aws.String("Salli"),                                                                   // Required
	}
	resp, err := svc.SynthesizeSpeech(params)

	if err != nil {
		fmt.Println(err.Error())
		return emptyBytes
	}
	data, err := ioutil.ReadAll(resp.AudioStream)

	if err != nil {
		fmt.Println(err.Error())
		return emptyBytes
	}
	cache.Set(cacheKey, data)
	return data
}

func formatLocation(name string) string {
	location := string(name[0:3])
	if location[2] == "-"[0] {
		location = name[0:4]
	}
	return strings.Replace(location, "-", "tac ", 1)
}

