package eintel

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Channel struct {
	Name      string
	directory string
	messages  chan ChannelMessage
	decoder   transform.Transformer
}

type ChannelMessage struct {
	Channel string
	Message string
}

var (
	pollInterval = 200 * time.Millisecond
)

func NewChannel(directory string, name string, messages chan ChannelMessage) *Channel {
	win16le := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewDecoder())

	return &Channel{
		Name:      name,
		directory: directory,
		messages:  messages,
		decoder:   utf16bom,
	}
}

func (c *Channel) Run() {

	fileChanged := make(chan os.FileInfo)
	monitorNewFile := make(chan os.FileInfo)
	fileContentChanged := make(chan int64)
	go c.monitorFileInfo(fileChanged)
	go c.monitorFileSize(monitorNewFile, fileContentChanged)
	var file os.FileInfo = nil
	for {
		select {
		case currentFile := <-fileChanged:
			file = currentFile
			monitorNewFile <- file
			fmt.Println("New file info for", c.Name)
			break
		case offset := <-fileContentChanged:
			if file.Size() > offset {
				c.ProcessChunk(file, offset)
			}
		}
	}

}

func (c *Channel) monitorFileSize(fileChanged chan os.FileInfo, newOffset chan int64) {

	var file os.FileInfo = nil
	oldSize := int64(0)

	for {
		if newFile, ok := <-fileChanged; ok {
			fmt.Println("File changed while monitoring size", c.Name)
			var size = int64(0)
			if file == nil {
				size = 0 //= newFile.Size()
			}
			file = newFile
			newOffset <- size
			oldSize = size
			continue
		}

		if file == nil {
			time.Sleep(pollInterval)
			continue
		}

		if stat, err := os.Stat(filepath.Join(c.directory, file.Name())); err == nil && oldSize != stat.Size() {
			newOffset <- oldSize
			oldSize = stat.Size()
		}
	}

}

func (c *Channel) monitorFileInfo(file chan os.FileInfo) {

	var currentFile os.FileInfo = nil
	var oldFile os.FileInfo = nil

	for {
		currentFile = c.latestFileRevision()
		if oldFile == nil || oldFile.Name() != currentFile.Name() {
			oldFile = currentFile
			file <- currentFile
		}
		time.Sleep(pollInterval)
	}

}

func (c *Channel) ProcessChunk(fileInfo os.FileInfo, offset int64) {
	fullName := filepath.Join(c.directory, fileInfo.Name())
	file, err := os.Open(fullName)
	defer file.Close()
	if err != nil {
		return
	}

	reader := bufio.NewReader(transform.NewReader(file, c.decoder))
	file.Seek(offset, io.SeekStart)

	text, err := ioutil.ReadAll(reader)

	for _, line := range strings.Split(string(text), `\r`) {
		c.messages <- ChannelMessage{
			Channel: c.Name,
			Message: line,
		}
	}

}

func (c *Channel) latestFileRevision() os.FileInfo {
	files, err := ioutil.ReadDir(c.directory)
	if err != nil {
		return nil
	}

	var latestFile os.FileInfo = nil
	latestVersion := uint64(0)
	for _, file := range files {
		path := filepath.Join(c.directory, file.Name())
		chanLogSpec := strings.TrimSuffix(file.Name(), filepath.Ext(path))
		tokens := strings.Split(chanLogSpec, "_")
		if len(tokens) != 3 || tokens[0] != c.Name {
			continue
		}

		version, err := strconv.ParseUint(fmt.Sprintf("%s%s", tokens[1], tokens[2]), 10, 64)
		if err != nil {
			continue
		}
		if latestVersion < version {
			latestVersion = version
			latestFile = file
		}
	}

	return latestFile
}
