package eintel

import (
	"fmt"
	"github.com/ogier/pflag"
	"github.com/vvlad/eintel/universe"
	"strings"
)

var (
	jumpAlert = 5
)

type IntelReport int

const (
	Unknown IntelReport = iota
	Alert
	Clear
)

type IntelMessage struct {
	Type        IntelReport
	Description string
	Location    *universe.System
}

type messageParser struct {
	Messages        chan IntelMessage
	input           chan ChannelMessage
	LocationChanges chan *universe.System
	currentLocation *universe.System
}

func NewMessageParser(chat *Chat, locations chan *universe.System) *messageParser {
	return &messageParser{
		input:           chat.Messages,
		Messages:        make(chan IntelMessage),
		LocationChanges: locations,
		currentLocation: universe.Systems["Jita"],
	}
}

func (m *messageParser) Run() {
	go m.locationMonitor()
	jumpCount := 0
	for channelMessage := range m.input {
		if system := m.findSystem(channelMessage.Message); system != nil {
			jumpCount = JumpCount(m.currentLocation.Name, system.Name)
			text := strings.ToLower(channelMessage.Message)
			if jumpCount > jumpAlert {
				continue
			}

			msg := IntelMessage{
				Location: system,
			}

			location := formatLocation(system.Name)
			if isClearReport(text) {
				msg.Type = Clear
				msg.Description = fmt.Sprintf("%s is clear", location)
			} else {
				msg.Type = Alert
				jumps := fmt.Sprintf("%d jumps", jumpCount)
				if jumpCount == 0 {
					jumps = fmt.Sprintf("1 jump")
				}
				msg.Description = fmt.Sprintf("Threat reported %s away in %s", jumps, location)
			}

			m.Messages <- msg
		}
	}
}

func formatLocation(name string) string {
	location := string(name[0:3])

	if location[2] == "-"[0] {
		location = name[0:4]
	}
	return strings.Replace(location, "-", "tac ", 1)
}

func isClearReport(text string) bool {
	words := strings.Fields(strings.ToLower(text))

	for _, word := range words {
		if word == "clr" || word == "clear" {
			return true
		}
	}

	return false
}

func (m *messageParser) locationMonitor() {
	for location := range m.LocationChanges {
		m.currentLocation = location
		fmt.Println("Channel changed to: ", location.Name)
	}
}

func (m *messageParser) findSystem(line string) *universe.System {

	for _, token := range strings.Fields(line) {
		if system, ok := universe.Systems[token]; ok {
			return system
		}
	}

	return nil
}

func init() {
	pflag.IntVar(&jumpAlert, "jumps", jumpAlert, "proximity of report")
}
