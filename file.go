package faye

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type block struct {
	start int64
	end   int64
	data  []byte
}

func (b *block) String() string {
	return fmt.Sprintf("bytes= %d-%d", b.start, b.end)
}

func (b *block) Parse(s string) error {
	tmp := strings.Split(s, "=")[1]
	tmp = strings.TrimSpace(tmp)
	locs := strings.Split(tmp, "-")
	var err error
	b.start, err = strconv.ParseInt(strings.TrimSpace(locs[0]), 10, 64)
	b.end, err = strconv.ParseInt(strings.TrimSpace(locs[1]), 10, 64)
	if err != nil {
		return err
	}
	return nil
}

func newBlock(s string) (*block, error) {
	b := &block{}
	err := b.Parse(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type fileManager struct {
	file        *os.File
	dataChan    <-chan *block
	finishChan  chan<- *block
	context     context.Context
	length      int64
	currentSize int64
}

func newManager(c context.Context, rawPath, name string, dataChan chan *block) (manager *fileManager, err error) {
	file, err := os.OpenFile(fileAddr(rawPath, name), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}
	manager = new(fileManager)
	manager.context = c
	manager.file = file
	manager.dataChan = dataChan
	var ok bool
	manager.finishChan, ok = c.Value("finish").(chan *block)
	if !ok {
		panic("Can not get length of the file from context!")
	}
	//var ok bool
	//manager.length, ok = c.Value("length").(int64)
	//if !ok {
	//	panic("Can not get length of the file from context!")
	//}
	return
}

func (m *fileManager) saveData(b *block) {
	m.file.Seek(b.start, 0)
	m.file.Write(b.data)
}

func (m *fileManager) start() {
	var b *block
	for {
		b = <-m.dataChan
		if b == nil {
			select {
			case <-m.context.Done():
				m.close()
				return
			default:
				panic("Receiving nil data from channel!")
			}
		}
		m.saveData(b)
		go func() {
			m.finishChan <- b
		}()
	}
}

func (m *fileManager) close() {
	m.file.Close()
}
