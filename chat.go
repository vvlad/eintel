package eintel

import (
	"errors"
	"github.com/mattes/go-expand-tilde"
	"github.com/fsnotify/fsnotify"
  "io/ioutil"
	"os"
	"time"
  "log"
  "path/filepath"
)

var (
	pollInterval   = 1 * time.Second
)

type PlayerChanels map[string]*Channel

type Chat struct {
	lines     chan string
	channels  []string
  knownChannels map[string]*Channel
	directory string
  localChannels map[string]*LocalChannel
  intelChannels map[string]*IntelChannel
  threatAssement *ThreatAssement
  intelMessages chan IntelMessage
  Locations chan LocationMessage
}

func NewChat(intel_messages chan IntelMessage) *Chat {

	directory, err := findChatLogsLocation()
	if err != nil {
		return nil
	}

	return &Chat{
		directory: directory,
		channels:  []string{},
    knownChannels: map[string]*Channel {},
    localChannels: map[string]*LocalChannel{},
    intelChannels: map[string]*IntelChannel{},
    intelMessages : intel_messages,
    Locations: make(chan LocationMessage),
	}

}

var (
	chatLogsLocations = []string{
		"~/Documents/EVE/logs/Chatlogs",
	}
)

func (c *Chat) AddChannel(name string) {
  c.channels = append(c.channels, name)
}

func (c *Chat) broadcastLocationChanges(local *LocalChannel, intel * IntelChannel) {
  for message := range local.Messages {
    intel.Locations <- message
    //c.Locations <- message
  }
}


func (c *Chat) broadcastIntelMessages(intel * IntelChannel) {
  for message := range intel.Messages {
  	c.intelMessages <- message
  }
}

func (c *Chat) Run() {

  for id, info := range playerChannelsWithName(c.directory, "Local") {
    parser := NewLocalChannel(info)
    intel := NewIntelChannel(info.PlayerName)

    go c.broadcastLocationChanges(parser, intel)
    go c.broadcastIntelMessages(intel)

    c.localChannels[info.PlayerName] = parser
    c.intelChannels[info.PlayerName] = intel
    c.knownChannels[id] = NewChannel(info, parser)
  }

  for _, name := range c.channels {
    channels := playerChannelsWithName(c.directory, name)
    for id, info := range channels {
      if intel, ok := c.intelChannels[info.PlayerName]; ok {
        player_channel := NewChannel(info, intel)
        c.knownChannels[id] = player_channel
      }
    }

  }

  for _, channel := range c.knownChannels {
    channel.Resume()
  }

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events: go c.distachToChannel(event)
			case err := <-watcher.Errors: log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(c.directory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func (c *Chat) distachToChannel(event fsnotify.Event) {
  info := ChannelInfoFromFile(event.Name)
  if info == nil { return }

  if channel, ok := c.knownChannels[info.Id] ; ok {
    channel.NotifyChanges(info)
  }
}

func playerChannelsWithName(directory, name string) (channels map[string]*ChannelInfo) {
  channels = map[string]*ChannelInfo{}
  files, err := ioutil.ReadDir(directory)
  if err != nil { return }
  for _, file := range files {
    info := ChannelInfoFromFile(filepath.Join(directory, file.Name()))
    if info == nil { continue }
    if info.Name != name { continue }
    if existing, ok := channels[info.Id]; ok {
      if existing.Version < info.Version { channels[info.Id] = info }
    } else {
      channels[info.Id] = info
    }
  }

  return
}

func findChatLogsLocation() (string, error) {
	for _, location := range chatLogsLocations {
		if path, err := tilde.Expand(location); err != nil {
			continue
		} else {
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}
	return "", errors.New("no chat location found")
}

