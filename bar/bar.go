package bar

import (
	"fmt"
	"log"
	"time"
)

type Bar struct {
	percent int64
	current int64
	total   int64
	rate    string
	graph   string
	start   time.Time
}

func (b *Bar) NewOption(start, total int64) {
	b.current = start
	b.total = total
	if b.graph == "" {
		b.graph = "█"
	}
	b.percent = b.getPercent()
	for i := 0; i < int(b.percent); i += 2 {
		b.rate += b.graph
	}
}

func (b *Bar) getPercent() int64 {
	return int64(float32(b.current) / float32(b.total) * 100)
}

func (b *Bar) NewOptionWithGraph(start, total int64, graph string) {
	b.graph = graph
	b.NewOption(start, total)
}

func (b *Bar) Start() {
	b.start = time.Now()
	log.Printf("Task starts at %v\n", b.start)
}

func (b *Bar) Play(cur int64) {
	b.current = cur
	last := b.percent
	b.percent = b.getPercent()
	if b.percent != last && b.percent%2 == 0 {
		b.rate += b.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", b.rate, b.percent, b.current, b.total)
}

func (b *Bar) Finish() {
	fmt.Println()
	log.Printf("Task finished take %v.\n", time.Since(b.start))
}
