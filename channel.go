package eintel

import (
	//  "log"
	"strings"
)

type WatchBehaviour int

const (
	Tail WatchBehaviour = iota
	Replay
)

type ChannelProcessor interface {
	UpdateInfo(info *ChannelInfo)
	ParseLine(line string)
}

type Channel struct {
	info   *ChannelInfo
	file   *ChannelFile
	parser ChannelProcessor
}

func NewChannel(info *ChannelInfo, parser ChannelProcessor) *Channel {
	parser.UpdateInfo(info)

	channel := &Channel{
		info:   info,
		file:   NewChannelFile(info),
		parser: parser,
	}
	return channel
}

func (c *Channel) Resume() {
	if c.info.Name != "Local" {
	 c.file.Resume()
	} else {
    c.NotifyChanges(c.info)
  }
}

func (c *Channel) NotifyChanges(info *ChannelInfo) {

	if c.info.Path != info.Path {
		c.info = info
		c.parser.UpdateInfo(info)
		c.file = NewChannelFile(info)
	}

	updates, err := c.file.ReadUpdates()
	if err != nil {
    log.Errorf("[%s] error %s while reading %s", info.PlayerName, err, info.Path)
		return
	}

	for {
		line, err := updates.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			c.parser.ParseLine(line)
		}
	}

}
