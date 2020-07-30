package faye

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestMaster(t *testing.T) {
	rawURL := "https://sel-kor-ping.vultr.com/vultr.com.100MB.bin"
	addr, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	master, err := NewMaster(rawURL, addr, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	master.Start()
}
