package main

import (
	"fmt"
	"net/http"

	"github.com/closetool/faye"
)

func main() {
	faye.Headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	faye.Headers.Add("referer", "https://www.bilibili.com/video/BV18K4y1x7ZM?spm_id_from=333.851.b_7265706f7274466972737431.8")
	faye.Headers.Add("origin", "https://www.bilibili.com")
	//Headers.Add("range", "bytes=1211770-1211779")
	rawURL := "https://upos-sz-mirrorcos.bilivideo.com/upgcxcode/60/80/216148060/216148060_nb2-1-30280.m4s?e=ig8euxZM2rNcNbdlhoNvNC8BqJIzNbfqXBvEqxTEto8BTrNvN0GvT90W5JZMkX_YN0MvXg8gNEV4NC8xNEV4N03eN0B5tZlqNxTEto8BTrNvNeZVuJ10Kj_g2UB02J0mN0B5tZlqNCNEto8BTrNvNC7MTX502C8f2jmMQJ6mqF2fka1mqx6gqj0eN0B599M=&uipk=5&nbs=1&deadline=1595646519&gen=playurl&os=cosbv&oi=1857883589&trid=a7f20d4d6274410cafeb061a0e61b0b7u&platform=pc&upsig=4b1c0e31497691df2795c5970abfb2fb&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,platform&mid=6846013&orderid=0,3&logo=80000000"
	thread := 8
	//addr := `C:\Users\Administration\Desktop`
	addr := `/mnt/c/Users/Administration/Desktop`
	client := &http.Client{}
	master, err := faye.NewMaster(rawURL, thread, addr, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	master.Start()
}
