package eintel

import (
	"bufio"
	// "fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type MessageProcessor interface {
	Process(c *Channel, line string)
}

type WatchBehaviour int

const (
	Tail WatchBehaviour = iota
	Replay
)

type Channel struct {
	Name      string
	chat      *Chat
	decoder   transform.Transformer
	file      *os.File
	fileName  string
	seekPos   int64
	Processor MessageProcessor
	Behaviour WatchBehaviour
}

type ChannelMessage struct {
	Channel string
	Message string
	TS      time.Time
	Sender  string
}

var (
	intelMessageFormat = regexp.MustCompile(`\[ (?P<date>\d+.\d+.\d+) (?P<time>\d+:\d+:\d+) \] (?P<name>.*?) > (?P<message>.*)`)
)

func NewChannel(chat *Chat, name string) *Channel {
	win16le := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewDecoder())

	return &Channel{
		Name:      name,
		chat:      chat,
		decoder:   utf16bom,
		fileName:  "",
		Processor: &intelProcessor{},
		Behaviour: Tail,
	}

}

func (c *Channel) Run() {

	for {
		if file, err := c.logFile(); err == nil {
			reader := bufio.NewReader(transform.NewReader(file, c.decoder))
			if data, err := ioutil.ReadAll(reader); err == nil && len(data) > 0 {
				for _, line := range strings.Split(string(data), "\n") {
					line = strings.TrimSpace(line)
					if len(line) > 0 {
						c.Processor.Process(c, line)
					}
				}
				if stat, err := file.Stat(); err == nil {
					c.seekPos = stat.Size()
				}
			}
		}
		time.Sleep(pollInterval)
	}
}

type intelProcessor struct {
}

func (i *intelProcessor) Process(c *Channel, line string) {
	matches := intelMessageFormat.FindStringSubmatch(line)

	if len(matches) > 0 {
		c.chat.Messages <- ChannelMessage{
			Channel: c.Name,
			Message: matches[4],
			Sender:  matches[3],
		}
	}

}

func (c *Channel) logFile() (*os.File, error) {

	if fileName, err := c.chat.findChannelLog(c.Name); err == nil {
		if info, err := os.Stat(fileName); err == nil {

			firstTime := c.file == nil

			if fileName != c.fileName {
				if c.file != nil {
					c.file.Close()
				}
				c.seekPos = 0

				if c.file, err = os.Open(fileName); err != nil {
					return nil, err
				}

				if firstTime && c.Behaviour == Tail {
					c.seekPos = info.Size()
				}
				c.fileName = fileName
			}
			c.file.Seek(c.seekPos, os.SEEK_SET)

		} else {
			return nil, err
		}

	} else {
		return nil, err
	}

	return c.file, nil

}
