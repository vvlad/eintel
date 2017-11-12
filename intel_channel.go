package eintel

import (
	"github.com/vvlad/eintel/universe"
	"regexp"
	"strings"
)

var (
	messagePattern = regexp.MustCompile(`\[ (?P<date>\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2}) \] (?P<player>.+) > (?P<message>.*)`)
)

type IntelChannel struct {
	PlayerName string
	Locations  chan LocationMessage
	System     *universe.System
	Messages   chan IntelMessage
}

type IntelMessage struct {
	PlayerName    string
	PlayerSystem  *universe.System
	RelatedSystem *universe.System
	Tokens        []string
	Line          string
}

func NewIntelChannel(playerName string) *IntelChannel {
	intel := &IntelChannel{
		PlayerName: playerName,
		Locations:  make(chan LocationMessage),
		System:     universe.Systems["Jita"],
		Messages:   make(chan IntelMessage),
	}
	go intel.Run()
	return intel
}

func (i *IntelChannel) Run() {
	for locationMessage := range i.Locations {
		i.System = locationMessage.System
		log.Noticef("Intel now knows that %s is in %s", i.PlayerName, i.System.Name)
	}
}

func (i *IntelChannel) UpdateInfo(info *ChannelInfo) {}

func (i *IntelChannel) ParseLine(line string) {
	if msg, ok := reSubMatchMap(messagePattern, line); ok {
		tokens := strings.Fields(msg["message"])
		tokens = Map(tokens, strings.ToUpper)
		tokens = Map(tokens, RemoveArtefacts)
		tokens = Filter(tokens, Without(universe.Ships...))
		tokens = Filter(tokens, Without(StopWords...))
		system := findSystem(tokens)
		if system == nil {
      log.Errorf("Unable to find system from %v", tokens)
			return
		}

		tokens = Filter(tokens, Without(strings.ToUpper(system.Name)))

		message := IntelMessage{
			PlayerName:    i.PlayerName,
			PlayerSystem:  i.System,
			RelatedSystem: system,
			Tokens:        tokens,
			Line:          line,
		}
		i.Messages <- message
    log.Debugf("[%s] intel message delivered", i.PlayerName)
	} else {
    log.Errorf("[%s] line %s doesn't match pattern", i.PlayerName, line)
  }
}

func RemoveArtefacts(word string) string {
	return strings.Replace(word, "*", "", -1)
}

func findSystem(tokens []string) *universe.System {
	for _, token := range tokens {
		token := strings.Replace(token, "*", "", -1)
		if system, ok := universe.Systems[token]; ok {
			return system
		}
	}

	return nil
}
