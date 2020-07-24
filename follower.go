package faye

import (
	"context"
	"net/http"
)

//Follower Denote the worker node which download the file block in a single thread
type Follower struct {
	context  context.Context
	client   *http.Client
	dataChan chan byte
}

func NewFollower(context context.Context, dataChan chan byte, client *http.Client) *Follower {
	fw := Follower{
		context:  context,
		client:   client,
		dataChan: dataChan,
	}
	return &fw
}
