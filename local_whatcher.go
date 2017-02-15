package eintel

import (
	"github.com/vvlad/eintel/universe"
	"regexp"
)

type localChannel struct {
	*Channel
	Updates chan *universe.System
}

var (
	localSystemMessageFormat = regexp.MustCompile(`EVE System > Channel changed to Local : (.*)`)
	localSystemHeaderFormat  = regexp.MustCompile(`Channel ID:      \(\('solarsystemid2', (\d+)\),\)`)
)

func NewLocalWatcher(chat *Chat) *localChannel {

	lch := &localChannel{
		Channel: NewChannel(chat, "Local"),
		Updates: make(chan *universe.System),
	}
	lch.Processor = lch
	lch.Behaviour = Replay

	return lch
}

func (c *localChannel) Process(channel *Channel, line string) {
	if localSystemHeaderFormat.MatchString(line) {
		c.Updates <- universe.Systems[localSystemHeaderFormat.FindStringSubmatch(line)[1]]
	}

	if localSystemMessageFormat.MatchString(line) {
		c.Updates <- universe.Systems[localSystemMessageFormat.FindStringSubmatch(line)[1]]
	}
}
