package eintel

import ()

type ThreatLevel int

const (
	ThreatLevelUnknown ThreatLevel = iota
  ThreatLevelIrelevant
	ThreatLevelCleared
	ThreatLevelDanger
	ThreatLevelWarning
	ThreatLevelFatal
)

var (
	ProximityWarning = 5
	ProximitDanger   = 3
)

type ThreatAssement struct {
	IntelMessages chan IntelMessage
  ThreatMessages chan ThreatMessage
}

type ThreatMessage struct {
	IntelMessage
	Level ThreatLevel
	Jumps int
}

func NewThreatAssement() *ThreatAssement {
	t := &ThreatAssement{
		IntelMessages: make(chan IntelMessage),
    ThreatMessages: make(chan ThreatMessage),
	}
	go t.run()
	return t
}

func (t *ThreatAssement) run() {
	for message := range t.IntelMessages {
    go t.assesThread(message)
	}
}

func (t *ThreatAssement) assesThread(message IntelMessage) {
  threat_level := ThreatLevelUnknown
  jumps := IrelevantNumberOfJumps

  if isIrelevant(message) {
    log.Debugf("Irelevant message %v", message)
    return
  }

  jumps = JumpCount(message.PlayerSystem.Name, message.RelatedSystem.Name)

  if jumps > ProximitDanger {
    return
  }

  threat_level = ThreatLevelIrelevant

  if isNoThreat(message) {
    threat_level = ThreatLevelCleared
  } else if jumps == UnknownNumberOfJumps {
    threat_level = ThreatLevelDanger
  } else if jumps <= ProximityWarning {
    threat_level = ThreatLevelWarning
  } else if jumps <= ProximitDanger {
    threat_level = ThreatLevelDanger
  } else if jumps == 0 {
    threat_level = ThreatLevelFatal
  }

  threat_message := ThreatMessage{
    IntelMessage: message,
    Jumps:        jumps,
    Level:        threat_level,
  }

  t.ThreatMessages <- threat_message
  log.Noticef("[%s] Threat assement delivered", message.PlayerName)
}

func isNoThreat(message IntelMessage) bool {
	return Intersects(message.Tokens, NoThreatList...)
}

func isIrelevant(message IntelMessage) bool {
	return Intersects(message.Tokens, IrelevantWordList...)
}

var IrelevantWordList = []string{
	"STATUS", "LOC", "LOCATION",
}

var NoThreatList = []string{
	"CLR", "CLEAR",
}
