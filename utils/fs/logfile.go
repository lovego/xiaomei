package fs

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type LogFile struct {
	path string
	sync.RWMutex
	*os.File
}

func NewLogFile(path string) (*LogFile, error) {
	l := LogFile{path: path}
	if err := l.open(); err != nil {
		return nil, err
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGUSR1)

	go func() {
		for {
			<-ch
			if err := l.open(); err != nil {
				log.Println(err)
			}
		}
	}()
	return &l, nil
}

func (l *LogFile) open() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0775); err != nil {
		return err
	}
	file, err := os.OpenFile(l.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	l.Lock()
	defer l.Unlock()
	if l.File != nil {
		if err := l.File.Close(); err != nil {
			log.Println(err)
		}
	}
	l.File = file
	return nil
}

func (l *LogFile) Write(b []byte) (n int, err error) {
	l.RLock()
	defer l.RUnlock()
	return l.File.Write(b)
}
