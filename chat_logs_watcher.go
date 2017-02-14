package eintel

import (
	// "bytes"
	"errors"
	// "fmt"
	"github.com/mattes/go-expand-tilde"
	// "golang.org/x/text/encoding/unicode"
	// "golang.org/x/text/transform"
	// "io/ioutil"
	"os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"
	"sync"
)

type ChatLogsWatcher struct {
	lines    chan string
	channels []*Channel
	Messages chan ChannelMessage
}

func NewChatLogsWatcher(channelNames ...string) *ChatLogsWatcher {

	directory, err := findChatLogsLocation()
	lines := make(chan ChannelMessage)

	if err != nil {
		return nil
	}

	channels := make([]*Channel, 0)
	for _, name := range channelNames {
		channel := NewChannel(directory, name, lines)
		channels = append(channels, channel)
	}

	return &ChatLogsWatcher{
		channels: channels,
		Messages: lines,
	}
}

var (
	chatLogsLocations = []string{
		"~/Documents/EVE/logs/Chatlogs",
	}
)

func (c *ChatLogsWatcher) Run() {
	var wg sync.WaitGroup

	for _, channel := range c.channels {
		go func(channel *Channel) {
			defer wg.Done()
			channel.Run()
		}(channel)
		wg.Add(1)
	}
	wg.Wait()
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
