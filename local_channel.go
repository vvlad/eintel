package eintel

import (
  "log"
	"github.com/vvlad/eintel/universe"
	"regexp"
)

type LocationMessage struct{
  Player string
  System *universe.System
}

type LocalChannel struct {
  info *ChannelInfo
  Messages chan LocationMessage
}

var (
	localSystemMessageFormat = regexp.MustCompile(`EVE System > Channel changed to Local : (?P<name>.*)`)
	localSystemHeaderFormat  = regexp.MustCompile(`\(\('solarsystemid2', (?P<id>\d+)\),\)`)
)

func NewLocalChannel(info *ChannelInfo) *LocalChannel {
  return &LocalChannel{
    info: info,
    Messages: make(chan LocationMessage),
  }
}

func (c *LocalChannel) UpdateInfo(info *ChannelInfo) {
  c.info = info
  if match, ok := reSubMatchMap(localSystemHeaderFormat, info.ChannelId); ok {
    c.broadCastSystemUpdate(info, universe.Systems[match["id"]])
  }
}

func (c *LocalChannel) ParseLine(line string) {
  if match, ok := reSubMatchMap(localSystemMessageFormat, line); ok {
    c.broadCastSystemUpdate(c.info, universe.Systems[match["name"]])
  }
}

func (c *LocalChannel) broadCastSystemUpdate(info *ChannelInfo, system *universe.System) {
  if info == nil || system == nil { return }
  log.Printf("%s is in %s", info.PlayerName, system.Name)
  c.Messages <- LocationMessage{
    Player: info.PlayerName,
    System: system,
  }
}
