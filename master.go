package faye

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	file     *os.File
	length   int64
	addr     string
	client   *http.Client
	dataChan chan byte
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
	length, canMul, err := checkHead(rawURL, client)
	if err != nil {
		return nil, err
	}
	c := context.Background()
	c, cancelFunc := context.WithCancel(c)
	master := Master{
		url:        u,
		context:    c,
		cancelFunc: cancelFunc,
		addr:       addr,
		client:     client,
		length:     length,
		dataChan:   make(chan byte, int(math.Log10(float64(length)))),
	}
	if !canMul {
		return nil, nil
	}
	file, err := os.OpenFile(fileAddr(addr, name), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}
	master.file = file
	fws := make([]*Follower, t)
	for i := range fws {
		fws[i] = NewFollower(master.context, master.dataChan, client)
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
	} else {
		name = url[slashIndex+1:]
	}
	dotIndex := strings.LastIndex(name, ".")
	if dotIndex != -1 {
		name = name[:dotIndex]
	}
	return
}

func checkHead(url string, client *http.Client) (length int64, canMul bool, err error) {
	length, canMul = -1, false
	req, err := http.NewRequest("HEADER", url, nil)
	if err != nil {
		return
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
	if resp.Header.Get("Accept-Ranges") != "" {
		canMul = true
	}
	return
}
