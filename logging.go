package eintel

import (
	"github.com/op/go-logging"
  "os"
)

var (
  Log = logging.MustGetLogger("eintel")
  log = Log
)


func init() {
  format := logging.MustStringFormatter(
    `%{color}%{time:15:04:05} %{shortfile} ▶%{color:reset} %{message}`,
  )
	backend := logging.NewLogBackend(os.Stderr, "", 0)
  logging.SetBackend(backend)
  logging.SetFormatter(format)
}
