package eintel

import (
	"errors"
	"fmt"
	"github.com/mattes/go-expand-tilde"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	pollInterval   = 1 * time.Second
	chatLogPattern = regexp.MustCompile(`^(?P<name>.*?)_(?P<date>\d+)_(?P<time>\d+)\.txt$`)
)

type Chat struct {
	lines     chan string
	channels  []*Channel
	Messages  chan ChannelMessage
	Files     []os.FileInfo
	directory string
}

func NewChat(channelNames ...string) *Chat {

	directory, err := findChatLogsLocation()
	if err != nil {
		return nil
	}

	clw := &Chat{
		Messages:  make(chan ChannelMessage),
		directory: directory,
		Files:     make([]os.FileInfo, 0),
		channels:  make([]*Channel, 0),
	}

	for _, name := range channelNames {
		channel := NewChannel(clw, name)
		clw.channels = append(clw.channels, channel)
	}

	return clw
}

var (
	chatLogsLocations = []string{
		// "~/Documents/EVE/logs/Chatlogs-debug/",
		"~/Documents/EVE/logs/Chatlogs",
	}
)

func (c *Chat) Run() {
	runChannel := func(channel *Channel) {
		channel.Run()
	}

	for _, channel := range c.channels {
		go runChannel(channel)
	}

	c.pullDirectory()
}

func (c *Chat) pullDirectory() {
	for {

		if files, err := ioutil.ReadDir(c.directory); err == nil {
			c.Files = files
		}

		time.Sleep(pollInterval * 20)
	}
}

func (c *Chat) findChannelLog(name string) (string, error) {
	var fileInfo os.FileInfo = nil
	version := uint64(0)
	for _, current := range c.Files {

		matches := chatLogPattern.FindAllStringSubmatch(current.Name(), -1)

		if len(matches) != 1 {
			continue
		}

		tokens := matches[0]
		if tokens[1] != name {
			continue
		}
		newVersion, err := strconv.ParseUint(fmt.Sprintf("%s%s", tokens[2], tokens[3]), 10, 64)

		if err != nil {
			continue
		}
		if version < newVersion {
			version = newVersion
			fileInfo = current
		}
	}

	if fileInfo != nil {
		return filepath.Join(c.directory, fileInfo.Name()), nil
	} else {
		return "", os.ErrNotExist
	}
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
