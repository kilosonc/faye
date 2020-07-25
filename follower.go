package faye

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Follower Denote the worker node which download the file block in a single thread
type Follower struct {
	name        string
	context     context.Context
	url         string
	client      *http.Client
	dataChan    chan<- *block
	taskRelease <-chan *block
}

//NewFollower create a follower
func NewFollower(c context.Context, dataChan chan *block, taskRelease chan *block, client *http.Client) *Follower {
	fw := Follower{
		context:     c,
		client:      client,
		dataChan:    dataChan,
		taskRelease: taskRelease,
	}
	var ok bool
	fw.url, ok = c.Value("url").(string)
	if !ok {
		panic("Get url from context failed!")
	}
	return &fw
}

func (f *Follower) start() {
	for {
		select {
		case task := <-f.taskRelease:
			//log.Printf("follower get a task\n")
			res, err := f.retry(task, f.download)
			if err != nil {
				panic(err)
			}
			f.dataChan <- res
		case <-f.context.Done():
			f.close()
			return
		}
	}
}

func (f *Follower) retry(b *block, fn func(*block) (*block, error)) (*block, error) {
	var err error
	for i := 0; i < RetryTimes; i++ {
		res, err := fn(b)
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (f *Follower) download(b *block) (*block, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", f.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", b.start, b.end))
	for k, v := range Headers {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}
	var resp *http.Response
	resp, err = f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("Status Code isn't 2XX: %d", resp.StatusCode)
	}

	return b, err
}

func (f *Follower) close() {
	//log.Printf("Follower %s exit\n", f.name)
}
