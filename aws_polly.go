package eintel

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io/ioutil"
)

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

func MD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
