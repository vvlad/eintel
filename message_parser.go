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

type IntelMessage struct {
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
			switch {
			case contains(text, "clr", "clear"):
				{
					msg.Description = fmt.Sprintf("%s is clear", location)

				}
			case jumpCount == 0:
				{
					msg.Description = "Thread in local. DOCK DOCK DOCK!!!"
				}
			case jumpCount == 1:
				{
					msg.Description = fmt.Sprintf("Thread one jump away in %s", location)
				}
			default:
				{
					msg.Description = fmt.Sprintf("Threat %d jumps away in %s", jumpCount, location)
				}
			}

			if contains(text, "dock", "dck", "docked") {
				msg.Description = fmt.Sprintf("Docked %s", msg.Description)
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
	return strings.Replace(location, "-", " tac ", 1)
}

func contains(text string, whats ...string) bool {
	words := strings.Fields(strings.ToLower(text))

	for _, word := range words {
		for _, what := range whats {
			if strings.ToLower(word) == strings.ToLower(what) {
				return true
			}
		}
	}

	return false
}

func (m *messageParser) locationMonitor() {
	for location := range m.LocationChanges {
		if location != nil {
			m.currentLocation = location
			fmt.Println("Channel changed to: ", location.Name)
		}
	}
}

func (m *messageParser) findSystem(line string) *universe.System {

	for _, token := range strings.Fields(line) {
		token := strings.Replace(token, "*", "", -1)
		if system, ok := universe.Systems[token]; ok {
			return system
		}
	}

	return nil
}

func init() {
	pflag.IntVar(&jumpAlert, "jumps", jumpAlert, "proximity of report")
}
