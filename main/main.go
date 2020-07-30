package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/closetool/faye"
)

const (
	_1MB = 1 << 20
)

var (
	ErrToFewArgs   = errors.New("There are at least two args!")
	ErrCanGetWD    = errors.New("Can not get current dir!")
	ErrParseCmdArg = errors.New("CMD args' format wrong!")
)

func main() {
	c := new(cmd)
	if err := c.Parse(); err != nil {
		log.Printf("An error occured: %v\n", err)
		return
	}
	faye.Thread = *c.t
	faye.BlockSize = _1MB * *c.b
	faye.RetryTimes = *c.r
	client := &http.Client{Timeout: time.Second * 10}
	master, err := faye.NewMaster(c.url, *c.addr, client)
	if err != nil {
		log.Printf("An error occured: %v\n", err)
		return
	}
	err = master.Start()
	if err != nil {
		log.Printf("An error occured: %v\n", err)
		return
	}
}

type cmd struct {
	t    *int
	b    *int64
	h    *string
	r    *int
	addr *string
	url  string
}

func (c *cmd) Parse() error {
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		return ErrToFewArgs
	}
	wd, err := os.Getwd()
	if err != nil {
		return ErrCanGetWD
	}
	c.t = flag.Int("t", 8, " -t <num> threads number downloading file")
	c.b = flag.Int64("b", 1, " -b <num> block size")
	c.h = flag.String("h", "", " -h headers eg: -h Content-Length:1024;")
	c.r = flag.Int("r", 3, " -r <num> retry times")
	c.addr = flag.String("a", wd, " -a <string> wherever you want to save the file")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.PrintDefaults()
		return ErrParseCmdArg
	}
	c.url = flag.Args()[0]
	return nil
}

func usage() {
	fmt.Printf("Usage: faye [OPTION] <URL>\n")
	fmt.Printf(" -a, --addr <string> where save the file")
	fmt.Printf(" -t, --thread <num> threads number that you want\n")
	fmt.Printf(" -b, --block <num> block size\n")
	fmt.Printf(" -h, --header <headers> headers in http request\n")
	fmt.Printf(" -r --retry <num> retry times\n")
}
