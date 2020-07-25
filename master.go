package faye

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrURLFormatError = errors.New("Url Format Wrong!")
	ErrSendReqFailed  = errors.New("Send Request Failed!")
)

//Master is job allocation node in main thread

//Master will create Followers under the control of flag -t
//and allocate the tasks to each follower
//master is also responsible for storing data
type Master struct {
	//Download target
	url *url.URL
	//context control followers
	context    context.Context
	cancelFunc context.CancelFunc
	followers  []*Follower
	//addresss that save the file
	fm          *fileManager
	client      *http.Client
	dataChan    chan *block
	finishChan  chan *block
	taskRelease chan *block
	blockCount  int
	blockTable  sync.Map
	length      int64
	canMul      bool
}

//NewMaster creates a new master node
func NewMaster(rawURL string, t int, addr string, client *http.Client) (*Master, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	name, err := assign(rawURL)
	if err != nil {
		return nil, err
	}
	//length, canMul, err := checkHead(rawURL, client)
	//if err != nil {
	//	return nil, err
	//}
	finishChan := make(chan *block, t)
	c := context.Background()
	c, cancelFunc := context.WithCancel(c)
	c = context.WithValue(c, "url", rawURL)
	c = context.WithValue(c, "finish", finishChan)
	master := Master{
		url:        u,
		context:    c,
		cancelFunc: cancelFunc,
		client:     client,
		//length:     length,
		dataChan:    make(chan *block, t),
		finishChan:  finishChan,
		taskRelease: make(chan *block, t),
	}
	//if !canMul {
	//	return nil, nil
	//}
	master.fm, err = newManager(master.context, addr, name, master.dataChan)
	if err != nil {
		return nil, err
	}
	fws := make([]*Follower, t)
	for i := range fws {
		fws[i] = NewFollower(master.context, master.dataChan, master.taskRelease, client)
	}
	master.followers = fws
	return &master, nil
}

func assign(url string) (string, error) {
	name, err := searchName(url)
	if err != nil {
		return "", err
	}
	return name, nil
}

func searchName(url string) (name string, err error) {
	slashIndex := strings.LastIndex(url, "/")
	if slashIndex == -1 {
		return "", ErrURLFormatError
	}
	name = url[slashIndex+1:]

	qIndex := strings.LastIndex(name, "?")
	if qIndex != -1 {
		name = name[:qIndex]
	}
	//dotIndex := strings.LastIndex(name, ".")
	//if dotIndex != -1 {
	//	name = name[:dotIndex]
	//}
	return
}

func checkHead(url string, client *http.Client) (length int64, canMul bool, err error) {
	length, canMul = -1, false
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	for k, v := range Headers {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		err = ErrSendReqFailed
		return
	}
	length, err = strconv.ParseInt(resp.Header.Get("Content-length"), 10, 64)
	if err != nil {
		return
	}
	//if resp.Header.Get("Accept-Ranges") != "" {
	canMul = true
	//}
	return
}

func (m *Master) allocate() {
	m.blockTable.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			panic("Get key from map failed!")
		}
		v, ok := value.(bool)
		if !ok {
			panic("Get value from map failed!")
		}
		if !v {
			b, err := newBlock(k)
			if err != nil {
				panic("Make block failed!")
			}
			go func() {
				m.taskRelease <- b
				//log.Printf("allocate task %s\n", b)
			}()
		}
		return true
	})
}

func (m *Master) init() error {
	m.blockCount = 0
	length, canMul, err := checkHead(m.url.String(), m.client)
	if err != nil {
		return err
	}
	if !canMul {
		panic("can't multiply thread!")
	}
	m.length = length
	m.fm.length = length
	m.canMul = canMul

	start := int64(0)
	end := int64(-1)
	for end < m.length-1 {
		start = end + 1
		end = start + BlockSize
		if end > m.length-1 {
			end = m.length - 1
		}
		m.blockTable.Store(fmt.Sprintf("bytes= %d-%d", start, end), false)
		m.blockCount++
	}
	return nil
}

func (m *Master) handleFinished() {
	finished := 0
	for {
		b := <-m.finishChan
		//log.Printf("block %s has completed\n", b)
		finished++
		m.blockTable.Store(b.String(), true)
		if finished == m.blockCount {
			//all tasks finished
			//log.Printf("All tasks are completed\n")
			m.close()
			return
		}
	}
}

func (m *Master) close() {
	close(m.dataChan)
	m.cancelFunc()
}

func (m *Master) Start() error {
	//log.Printf("Master start\n")
	err := m.init()
	if err != nil {
		return err
	}
	//log.Printf("Init succeed")
	m.allocate()
	//log.Printf("allocate complete")
	for _, f := range m.followers {
		go f.start()
	}
	go m.fm.start()
	m.handleFinished()
	return nil
}
