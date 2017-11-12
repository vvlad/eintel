package eintel

import (
	"log"
)

type ThreatLevel int

const (
  ThreatLevelUnknown ThreatLevel = iota
  ThreatLevelDanger
  ThreatLevelWarning
  ThreatLevelFatal
  ThreatLevelCleared
)

var (
	ProximityWarning = 5
	ProximitDanger = 3
)

type ThreatAssement struct{
  IntelMessages chan IntelMessage
}

type ThreatMessage struct{
  IntelMessage
  Level ThreatLevel
  Jumps int
}

func NewThreatAssement() *ThreatAssement {
  t := &ThreatAssement{
		IntelMessages: make(chan IntelMessage),
	}
  go t.run()
  return t
}

func (t* ThreatAssement) run() {
  for message := range t.IntelMessages {
    threat_level := ThreatLevelUnknown
    jumps := IrelevantNumberOfJumps

    if isIrelevant(message) { return }

    if isNoThreat(message) {
      threat_level = ThreatLevelCleared
    }else{

      jumps = JumpCount(message.PlayerSystem.Name, message.RelatedSystem.Name)

      if jumps == UnknownNumberOfJumps { threat_level = ThreatLevelDanger }
      if jumps <= ProximityWarning { threat_level = ThreatLevelWarning }
      if jumps <= ProximitDanger { threat_level = ThreatLevelDanger }
      if jumps == 0 { threat_level = ThreatLevelFatal }
    }

    message := ThreatMessage{
      IntelMessage: message,
      Jumps: jumps,
      Level: threat_level,
    }

		log.Printf("Thread leve is %d jumps %d for %s %v", threat_level, jumps, message.Line, message.PlayerName)
  }
}

func isNoThreat(message IntelMessage) bool {
  return Intersects(message.Tokens, NoThreadList...)
}

func isIrelevant(message IntelMessage) bool {
  return Intersects(message.Tokens, IrelevantWordList...)
}


var IrelevantWordList = []string{
  "STATUS",
}

var NoThreadList = []string{
  "CLR", "CLEAR",
}


