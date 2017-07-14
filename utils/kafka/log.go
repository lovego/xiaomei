package kafka

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/lovego/xiaomei/utils/fs"
)

func (c *Consume) setupLogFile() {
	if dir := filepath.Dir(c.LogPath); dir != `.` {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
	if logFile, err := fs.OpenAppend(c.LogPath); err != nil {
		log.Fatal(err)
	} else {
		c.logFile = logFile
	}
}

func (c *Consume) writeLog(m map[string]interface{}) {
	buf, err := json.Marshal(m)
	if err != nil {
		log.Printf("marshal log err: %v", err)
		return
	}
	buf = append(buf, '\n')
	_, err = c.logFile.Write(buf)
	if err != nil {
		log.Printf("write log err: %v", err)
		return
	}
}
