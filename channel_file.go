package eintel

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ChannelInfo struct {
	Id             string
	Name           string
	PlayerName     string
	SessionStarted string
	Path           string
	ChannelId      string
	Offset         int64
	Version        int64
}

type ChannelFile struct {
	info   *ChannelInfo
	Offset int64
}

var (
	headerPattern = regexp.MustCompile(`(?sm)
\s+Channel ID:\s+(?P<id>.+)
\s+Channel Name:\s+(?P<name>.+)
\s+Listener:\s+(?P<player>.+)
\s+Session started:\s+(?P<date>\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})
`)

	versionPattern = regexp.MustCompile(`^(?P<name>.*?)_(?P<date>\d+)_(?P<time>\d+)\.txt$`)
)

func NewChannelFile(info *ChannelInfo) *ChannelFile {
	return &ChannelFile{info, info.Offset}
}

func (c *ChannelFile) Resume() {

	file, err := os.Open(c.info.Path)
	if err != nil {
		return
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		c.Offset = stat.Size()
	}

}

func (c *ChannelFile) ReadUpdates() (bytes.Buffer, error) {

	var buffer bytes.Buffer

	file, err := os.Open(c.info.Path)
	if err != nil {
		return buffer, err
	}

	file.Seek(c.Offset, io.SeekCurrent)
	defer file.Close()

	reader := NewUTF16Reader(file)
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}

	if stat, err := file.Stat(); err == nil {
		c.Offset = stat.Size()
	}
	return buffer, nil
}

func ChannelInfoFromFile(path string) *ChannelInfo {
	file, _ := os.Open(path)
	var info *ChannelInfo

	var buffer bytes.Buffer
	reader := NewUTF16Reader(file)
	reading_header := false
	offset := int64(0)

	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 3 && line[:3] == "---" {
			if reading_header {
				offset, err = file.Seek(0, io.SeekCurrent)
				break
			}
			reading_header = true
			continue
		}

		if reading_header {
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
	}

	if result, ok := reSubMatchMap(headerPattern, buffer.String()); ok {
		var id string
		id = fmt.Sprintf("%s-%s", result["name"], result["player"])
		version := int64(0)

		if v, ok := reSubMatchMap(versionPattern, path); ok {
			version, _ = strconv.ParseInt(fmt.Sprintf("%s%s", v["date"], v["time"]), 10, 64)
		}

		info = &ChannelInfo{
			Name:           result["name"],
			ChannelId:      result["id"],
			PlayerName:     result["player"],
			SessionStarted: result["date"],
			Path:           path,
			Id:             id,
			Offset:         offset,
			Version:        version,
		}
	}
	return info
}

func reSubMatchMap(r *regexp.Regexp, str string) (map[string]string, bool) {
	match := r.FindStringSubmatch(str)
	if len(match) < len(r.SubexpNames()) {
		return nil, false
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap, true
}
